package search

import (
	"math"

	"github.com/Droanox/ChessEngineAI/src/board"
)

func alphabeta(alpha int, beta int, depth int, cb *board.ChessBoard) int {
	// pvLength[board.Ply+1] is used to store the length of the principal variation
	pvLength[board.Ply+1] = board.Ply + 1

	// when the depth is 0, we call quiescence search to search for captures
	if depth == 0 {
		return quiescence(alpha, beta, cb)
	}

	// If the maximum ply is reached, we evaluate the position and return the score
	// instead of searching further as it'll break the search

	nodes++

	// check if the side to move is in check
	// if so, we increase the depth by 1 to search deeper
	var isChecked bool
	if board.SideToMove == board.White {
		isChecked = cb.IsSquareAttackedBySide(board.BitScanForward(cb.WhiteKing), board.Black)
	} else {
		isChecked = cb.IsSquareAttackedBySide(board.BitScanForward(cb.BlackKing), board.White)
	}
	if isChecked {
		depth++
	}

	// legalMovesNum is used to count the number of legal moves
	var legalMovesNum int

	// movelist is used to store the list of moves generated by the GenerateMoves function
	var moveList = []board.Move{}
	// generate all the moves
	cb.GenerateMoves(&moveList)

	// if the principal variation was found in the previous iteration, we score the principal variation
	// otherwise, we score all the moves
	if pvFollowed {
		scorePV(&moveList)
	}
	scoreMoves(&moveList)

	// movesSearched is used to count the number of moves searched
	var movesSearched int

	// search through the moves
	for i := 0; i < len(moveList); i++ {
		pickMove(&moveList, i)

		if !cb.MakeMove(moveList[i]) {
			continue
		}

		// increment legal moves for checkmate and stalemate detection
		legalMovesNum++

		// store the score
		var score int

		// Principal Variation Search (PVS) and Late Move Reduction (LMR)
		// https://www.chessprogramming.org/Principal_Variation_Search
		// https://www.chessprogramming.org/Late_Move_Reductions
		// full depth search
		if movesSearched == 0 {
			score = -alphabeta(-beta, -alpha, depth-1, cb)
			// LMR
		} else {
			// if the move satisfies the LMR conditions, we search deeper
			if (movesSearched >= fullDepthMoves) && (depth >= reductionLimit) && ((moveList[i].GetMoveFlags() & (board.MoveCaptures | board.MoveKnightPromotion)) == 0) {
				score = -alphabeta(-alpha-1, -alpha, depth-2, cb)
			} else {
				score = alpha + 1
			}
			// if the move fails high, we search deeper
			// principal variation search (PVS)
			if score > alpha {
				score = -alphabeta(-alpha-1, -alpha, depth-1, cb)
				// if the move fails high, we search deeper again to confirm the move is good and not a fluke
				if (score > alpha) && (score < beta) {
					score = -alphabeta(-beta, -alpha, depth-1, cb)
				}
			}
		}

		// unmake the move
		cb.MakeBoard()

		// fails high
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

		movesSearched++
	}

	// check for checkmate and stalemate
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