package engine

import (
	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/eval"
)

func quiescence(alpha int, beta int, cb *board.ChessBoard) int {
	// check if the search should be stopped, time is checked every 10240 nodes
	nodes++

	var standPat int = eval.Eval(*cb)
	// fails high
	if standPat >= beta {
		return beta
	}
	// found a better move
	if standPat > alpha {
		alpha = standPat
	}

	var moveList = []board.Move{}
	cb.GenerateCaptures(&moveList)

	scoreMoves(&moveList, board.Move{})

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

		// found a better move
		if score > alpha {
			alpha = score

			// fails high
			if score >= beta {
				return beta
			}
		}
	}
	// fails low
	return alpha
}

func ZWSearch(beta int, depth int, cb *board.ChessBoard) int {
	if depth == 0 {
		return quiescence(beta-1, beta, cb)
	}

	var moveList = []board.Move{}
	cb.GenerateMoves(&moveList)

	for i := 0; i < len(moveList); i++ {
		if !cb.MakeMove(moveList[i]) {
			continue
		}

		var score int = -ZWSearch(1-beta, depth-1, cb)
		cb.MakeBoard()

		// check if the search should be stopped, time is checked concurrently
		if IsStopped {
			return 0
		}

		// found a better move
		if score >= beta {
			return beta
		}
	}
	return beta - 1
}
