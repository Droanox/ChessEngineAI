package search

import (
	"fmt"
	"math"

	"github.com/Droanox/ChessEngineAI/src/board"
)

func Search(depth int, cb *board.ChessBoard) {
	// reset nodes counter
	nodes = 0

	// reset killer moves and history moves
	killerMoves = [2][board.MaxPly]board.Move{}
	historyMoves = [12][board.MaxPly]int{}

	// reset principal variation
	pvTable = [board.MaxPly][board.MaxPly]board.Move{}
	pvLength = [board.MaxPly]int{}
	pvFollowed = false

	for currDepth := 1; currDepth <= depth; currDepth++ {
		// follow principal variation
		pvFollowed = true

		// perform negamax search
		var score int = negamax(math.MinInt32, math.MaxInt32, currDepth, cb)

		// print principal variation
		fmt.Printf("info depth %d nodes %d score cp %d pv ", currDepth, nodes, score)
		for i := 0; i < pvLength[board.Ply+1]; i++ {
			PrintMove(pvTable[0][i])
			fmt.Print(" ")
		}
		fmt.Println()
	}

	// print best move
	fmt.Printf("bestmove ")
	PrintMove(pvTable[0][0])
	fmt.Println()
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
