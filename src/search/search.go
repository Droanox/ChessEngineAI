package search

import (
	"fmt"

	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/engine"
)

func PrintMove(move board.Move) {
	var pieceToString = []string{
		board.Knight: "n", board.Bishop: "b", board.Rook: "r", board.Queen: "q",
	}

	if (move.GetMoveFlags() & 0b1000) > 0 {
		fmt.Printf("%s%s%s",
			board.IndexToSquare[move.GetMoveStart()],
			board.IndexToSquare[move.GetMoveEnd()],
			pieceToString[board.PromotionToPiece[move.GetMoveFlags()]])
	}
	fmt.Printf("%s%s",
		board.IndexToSquare[move.GetMoveStart()],
		board.IndexToSquare[move.GetMoveEnd()])
}

func Search(depth int, cb *board.ChessBoard) {
	_ = negamax(-100000, 100000, depth, cb)
	fmt.Printf("bestmove ")
	PrintMove(bestMove)
	fmt.Println()
}

func negamax(alpha int, beta int, depth int, cb *board.ChessBoard) int {
	if depth == 0 {
		return engine.Eval(*cb)
	}
	nodes++
	var oldAlpha int = alpha
	//fmt.Println(oldAlpha)
	var bestMoveSoFar board.Move
	// fmt.Println(bestMove)

	var moveList = []board.Move{}
	cb.GenerateMoves(&moveList)

	for i := 0; i < len(moveList); i++ {
		if !cb.MakeMove(moveList[i]) {
			continue
		}
		var score int = -negamax(-beta, -alpha, depth-1, cb)
		cb.MakeBoard()

		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score

			if board.Ply == -1 {
				bestMoveSoFar = moveList[i]
			}
		}
	}

	if oldAlpha != alpha {
		bestMove = bestMoveSoFar
	}

	return alpha
}
