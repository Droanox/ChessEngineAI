package engine

import (
	"fmt"
	"time"

	"github.com/Droanox/ChessEngineAI/src/board"
)

func Search(depth int, cb *board.ChessBoard) {
	// reset nodes counter
	nodes = 0

	// set max depth
	if depth > board.MaxPly-4 {
		depth = board.MaxPly - 4
	}

	// reset isStopped flag
	IsStopped = false

	// reset killer moves and history moves
	killerMoves = [2][board.MaxPly]board.Move{}
	// counterMoves = [2][64][64]board.Move{}
	hhScore = [64][64]int{}
	bfScore = [64][64]int{}

	// reset principal variation
	pvTable = [board.MaxPly][board.MaxPly]board.Move{}
	pvLength = [board.MaxPly]int{}

	// start timer
	start = time.Now()

	// set alpha and beta
	alpha := minScore // -INFINITY
	beta := maxScore  // INFINITY

	// start listening for stop signal
	isTimerOn := make(chan bool, 1)
	go listenForStop(isTimerOn)

	for currDepth := 1; currDepth <= depth; currDepth++ {
		// pvFollowed = true

		// perform negamax search
		var score int = alphabeta(alpha, beta, currDepth, StandardSearch, cb)

		if IsStopped {
			break
		}

		// aspiration window, if the score is outside the window, we search again
		if (score <= alpha) || (score >= beta) {
			alpha = minScore // -INFINITY
			beta = maxScore  // INFINITY
			currDepth--
			continue
		}

		// set window up for next iteration
		alpha = score - aspirationWindow
		beta = score + aspirationWindow

		if pvLength[0] > 0 {
			// print principal variation
			if score > -MateValue && score < -MateScore {
				fmt.Printf("info depth %d nodes %d score mate %d time %d pv ", currDepth, nodes, -(MateValue+score)/2-1, time.Since(start).Milliseconds())
			} else if score < MateValue && score > MateScore {
				fmt.Printf("info depth %d nodes %d score mate %d time %d pv ", currDepth, nodes, (MateValue-score)/2+1, time.Since(start).Milliseconds())
			} else {
				fmt.Printf("info depth %d nodes %d score cp %d time %d pv ", currDepth, nodes, score, time.Since(start).Milliseconds())
			}
			for i := 0; i < pvLength[0]; i++ {
				PrintMove(pvTable[0][i])
				fmt.Print(" ")
			}
			fmt.Println()
		}
	}
	// stop listening for stop signal
	isTimerOn <- true

	// print best move
	fmt.Printf("bestmove ")
	PrintMove(pvTable[0][0])
	fmt.Println()

	// fmt.Println(board.Ply, " ", board.RepetitionTableIndexOffset, time.Since(start).Milliseconds())
}

func listenForStop(ch chan bool) {
	for {
		select {
		case <-ch:
			return
		default:
			if TimeControl && (time.Since(start).Milliseconds() > StopTime) {
				IsStopped = true
			}
			// ease up on the CPU
			time.Sleep(200 * time.Microsecond)
		}
	}
}

func PrintMove(move board.Move) {
	fmt.Printf("%s%s",
		board.IndexToSquare[move.GetMoveStart()],
		board.IndexToSquare[move.GetMoveEnd()])
	if (move.GetMoveFlags() & 0b1000) > 0 {
		var pieceToString = []string{
			board.Knight: "n", board.Bishop: "b", board.Rook: "r", board.Queen: "q",
		}
		if board.SideToMove == board.White {
			fmt.Printf("%s",
				pieceToString[board.PromotionToPiece[move.GetMoveFlags()]])
		} else {
			fmt.Printf("%s",
				pieceToString[board.PromotionToPiece[move.GetMoveFlags()]])
		}
	}
}
