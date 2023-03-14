package engine

import (
	"github.com/Droanox/ChessEngineAI/src/board"
)

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

func scoreMoves(movelist *[]board.Move) {
	moves := (*movelist)

	for i := 0; i < len(moves); i++ {
		if moves[i].GetMoveCapturedPiece() != board.EmptyPiece {
			moves[i].Score = moveOrderOffset + MVV_LVA[moves[i].GetMoveCapturedPiece()][moves[i].GetMoveStartPiece()]
		} else {
			if moves[i] == killerMoves[0][board.Ply+1] {
				moves[i].Score = moveOrderOffset - 10
			} else if moves[i] == killerMoves[1][board.Ply+1] {
				moves[i].Score = moveOrderOffset - 20
			} else {
				moves[i].Score = historyMoves[moves[i].GetMoveStartPiece()+(6*board.SideToMove)-1][moves[i].GetMoveEnd()]
			}
		}
	}
}

func pickMove(movelist *[]board.Move, currentIndex int) {
	for nextMove := currentIndex + 1; nextMove < len(*movelist); nextMove++ {
		if (*movelist)[nextMove].Score > (*movelist)[currentIndex].Score {
			// swap elements in movelist
			(*movelist)[nextMove], (*movelist)[currentIndex] = (*movelist)[currentIndex], (*movelist)[nextMove]
		}
	}
}
