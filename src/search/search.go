package search

import (
	"fmt"
	"math"

	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/engine"
)

func Search(depth int, cb *board.ChessBoard) {
	nodes = 0
	var score int = negamax(math.MinInt32, math.MaxInt32, depth, cb)
	fmt.Printf("info depth %d nodes %d score cp %d\n", depth, nodes, score)

	fmt.Printf("bestmove ")
	PrintMove(bestMove)
	fmt.Println()
}

func negamax(alpha int, beta int, depth int, cb *board.ChessBoard) int {
	if depth == 0 {
		return quiescence(alpha, beta, cb)
	}
	nodes++

	var isChecked bool
	if board.SideToMove == board.White {
		isChecked = cb.IsSquareAttackedBySide(board.BitScanForward(cb.WhiteKing), board.Black)
	} else {
		isChecked = cb.IsSquareAttackedBySide(board.BitScanForward(cb.BlackKing), board.White)
	}
	if isChecked {
		depth++
	}

	var oldAlpha int = alpha
	var bestMoveSoFar board.Move
	var legalMovesNum int

	var moveList = []board.Move{}
	cb.GenerateMoves(&moveList)

	scoreMoves(&moveList)

	for i := 0; i < len(moveList); i++ {
		pickMove(&moveList, i)

		if !cb.MakeMove(moveList[i]) {
			continue
		}
		legalMovesNum++
		var score int = -negamax(-beta, -alpha, depth-1, cb)
		cb.MakeBoard()

		// fail-hard
		if score >= beta {
			if (moveList[i].GetMoveFlags() & board.MoveCaptures) == 0 {
				killerMoves[1][board.Ply+1] = killerMoves[0][board.Ply+1]
				killerMoves[0][board.Ply+1] = moveList[i]
			}

			return beta
		}
		// found a better move
		if score > alpha {
			if (moveList[i].GetMoveFlags() & board.MoveCaptures) == 0 {
				historyMoves[moveList[i].GetMoveStartPiece()+(6*board.SideToMove)-1][moveList[i].GetMoveEnd()] += depth
			}

			alpha = score

			if board.Ply == -1 {
				bestMoveSoFar = moveList[i]
			}
		}
	}

	if legalMovesNum == 0 {
		if isChecked {
			return (math.MinInt32 + 10) + board.Ply
		} else {
			return 0
		}
	}
	if oldAlpha != alpha {
		bestMove = bestMoveSoFar
	}

	// fails low
	return alpha
}

func quiescence(alpha int, beta int, cb *board.ChessBoard) int {
	nodes++
	var standPat int = engine.Eval(*cb)
	// fail-hard
	if standPat >= beta {
		return beta
	}
	// found a better move
	if alpha < standPat {
		alpha = standPat
	}

	var moveList = []board.Move{}
	cb.GenerateCaptures(&moveList)

	scoreMoves(&moveList)

	for i := 0; i < len(moveList); i++ {
		pickMove(&moveList, i)

		if !cb.MakeCapture(moveList[i]) {
			continue
		}
		var score int = -quiescence(-beta, -alpha, cb)
		cb.MakeBoard()

		// fail-hard
		if score >= beta {
			return beta
		}
		// found a better move
		if score > alpha {
			alpha = score
		}
	}
	// fails low
	return alpha
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
