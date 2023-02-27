package search

import (
	"math"

	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/engine"
)

func negamax(alpha int, beta int, depth int, cb *board.ChessBoard) int {
	pvLength[board.Ply+1] = board.Ply + 1

	if depth == 0 {
		return quiescence(alpha, beta, cb)
	}
	if board.Ply > board.MaxPly {
		return engine.Eval(*cb)
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

	var legalMovesNum int

	var moveList = []board.Move{}
	cb.GenerateMoves(&moveList)

	if pvFollowed {
		scorePV(&moveList)
	} else {
		scoreMoves(&moveList)
	}

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

			pvTable[board.Ply+1][board.Ply+1] = moveList[i]

			for j := board.Ply + 2; j < pvLength[board.Ply+2]; j++ {
				pvTable[board.Ply+1][j] = pvTable[board.Ply+2][j]
			}

			pvLength[board.Ply+1] = pvLength[board.Ply+2]
		}
	}

	if legalMovesNum == 0 {
		if isChecked {
			return (math.MinInt32 + 10) + board.Ply
		} else {
			return 0
		}
	}

	// fails low
	return alpha
}
