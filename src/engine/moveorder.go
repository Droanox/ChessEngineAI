package engine

import (
	"github.com/Droanox/ChessEngineAI/src/board"
)

/*
	func scorePV(movelist *[]board.Move) {
		pvFollowed = false
		moves := (*movelist)

		for i := 0; i < len(moves); i++ {
			if moves[i].Move == pvTable[0][board.Ply+1].Move {
				moves[i].Score = moveOrderOffset + 100
				pvFollowed = true
			}
		}
	}
*/
func scoreMoves(movelist *[]board.Move) {
	moves := (*movelist)

	if pvFollowed {
		pvFollowed = false
		for i := 0; i < len(moves); i++ {
			if moves[i].Move == pvTable[0][board.Ply].Move {
				moves[i].Score = moveOrderOffset + 1000
				pvFollowed = true
				return
			}
		}
	}

	for i := 0; i < len(moves); i++ {
		if moves[i].GetMoveCapturedPiece() != board.EmptyPiece {
			moves[i].Score = moveOrderOffset + MVV_LVA[moves[i].GetMoveCapturedPiece()][moves[i].GetMoveStartPiece()]
		} else {
			if moves[i].Move == killerMoves[0][board.Ply].Move {
				moves[i].Score = moveOrderOffset - 10
			} else if moves[i].Move == killerMoves[1][board.Ply].Move {
				moves[i].Score = moveOrderOffset - 20
			} else {
				moves[i].Score = historyMoves[moves[i].GetMoveStartPiece()+(6*board.SideToMove)-1][moves[i].GetMoveEnd()]
			}
		}
	}
}

func pickMove(movelist *[]board.Move, currentIndex int) {
	moves := (*movelist)

	// Move through the list and if there's a move greater than the current Index
	// swap the move, keep doing this until a cut off occurs
	for nextMove := currentIndex + 1; nextMove < len(moves); nextMove++ {
		if (moves)[nextMove].Score > (moves)[currentIndex].Score {
			// swap elements in movelist
			(moves)[nextMove], (moves)[currentIndex] = (moves)[currentIndex], (moves)[nextMove]
		}
	}
}
