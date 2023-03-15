package engine

import (
	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/eval"
)

func quiescence(alpha int, beta int, cb *board.ChessBoard) int {
	// check if the search should be stopped, time is checked every 10240 nodes
	/*
		if nodes&10240 == 0 {
			listenForStop()
		}
	*/
	nodes++
	var standPat int = eval.Eval(*cb)
	// fails high
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

		// check if the search should be stopped, time is checked concurrently
		if IsStopped {
			return 0
		}

		// fails high
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
