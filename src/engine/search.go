package engine

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/Droanox/ChessEngineAI/src/board"
)

func Search(depth int, cb *board.ChessBoard) {
	// reset nodes counter
	nodes = 0

	// reset isStopped flag
	isStopped = false

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

	// start listening for stop command
	isFinished := make(chan bool, 1)
	go parsecmd(isFinished)

	for currDepth := 1; currDepth <= depth; currDepth++ {
		// follow principal variation
		pvFollowed = true

		// perform negamax search
		var score int = alphabeta(alpha, beta, currDepth, cb)

		if isStopped {
			break
		}

		// aspiration window, if the score is outside the window, we search again
		if (score <= alpha) || (score >= beta) {
			alpha = math.MinInt32 // -INFINITY
			beta = math.MaxInt32  // INFINITY
			continue
		}

		// set window up for next iteration
		alpha = score - 50
		beta = score + 50

		// print principal variation
		fmt.Printf("info depth %d nodes %d score cp %d time %d pv ", currDepth, nodes, score, time.Since(start).Milliseconds())
		for i := 0; i < pvLength[board.Ply+1]; i++ {
			PrintMove(pvTable[0][i])
			fmt.Print(" ")
		}
		fmt.Println()
	}
	// stop listening for stop command
	isFinished <- true

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
				isStopped = true
			}
			// ease up on the CPU
			time.Sleep(1 * time.Millisecond)
		}
	}
}

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
				isStopped = true
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
