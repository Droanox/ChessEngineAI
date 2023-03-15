package engine

import (
	"fmt"
	"math"
	"time"

	"github.com/Droanox/ChessEngineAI/src/board"
)

func Init() {
	initHash()
}

func Search(depth int, cb *board.ChessBoard) {
	// reset nodes counter
	nodes = 0

	// reset isStopped flag
	IsStopped = false

	// reset killer moves and history moves
	killerMoves = [2][board.MaxPly]board.Move{}
	historyMoves = [12][board.MaxPly]int{}

	// reset principal variation
	pvTable = [board.MaxPly][board.MaxPly]board.Move{}
	pvLength = [board.MaxPly]int{}
	pvFollowed = false

	// start timer
	start = time.Now()

	// set alpha and beta
	alpha := math.MinInt32 // -INFINITY
	beta := math.MaxInt32  // INFINITY

	// start listening for stop signal
	isTimerOn := make(chan bool, 1)
	go listenForStop(isTimerOn)

	for currDepth := 1; currDepth <= depth; currDepth++ {
		// follow principal variation
		pvFollowed = true

		// perform negamax search
		var score int = alphabeta(alpha, beta, currDepth, cb)

		// check if the search should be stopped, time is checked concurrently
		if IsStopped {
			break
		}

		// aspiration window, if the score is outside the window, we search again
		if (score <= alpha) || (score >= beta) {
			alpha = math.MinInt32 // -INFINITY
			beta = math.MaxInt32  // INFINITY
			currDepth--
			continue
		}

		// set window up for next iteration
		alpha = score - aspirationWindow
		beta = score + aspirationWindow

		// print principal variation
		fmt.Printf("info depth %d nodes %d score cp %d time %d pv ", currDepth, nodes, score, time.Since(start).Milliseconds())
		for i := 0; i < pvLength[board.Ply+1]; i++ {
			PrintMove(pvTable[0][i])
			fmt.Print(" ")
		}
		fmt.Println()
	}
	// stop listening for stop signal
	isTimerOn <- true

	// print best move
	fmt.Printf("bestmove ")
	PrintMove(pvTable[0][0])
	fmt.Println()
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
			time.Sleep(1 * time.Millisecond)
		}
	}
}

/*
func parsecmd(ch chan bool) {
	exit := make(chan bool, 1)
	go listenForStop(exit)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-ch:
			exit <- true
			return
		default:
			if scanner.Text() == "stop" {
				IsStopped = true
				break
			}
			if scanner.Text() == "quit" {
				os.Exit(0)
			}
			// ease up on the CPU
			time.Sleep(1 * time.Millisecond)
		}
	}
}
*/

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
