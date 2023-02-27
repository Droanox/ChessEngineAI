package search

import (
	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/engine"
)

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
