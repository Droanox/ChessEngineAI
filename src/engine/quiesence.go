package engine

import (
	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/eval"
)

func quiescence(alpha int, beta int, cb *board.ChessBoard) int {
	// apply 3 move repetition rule and 50 move rule
	if board.Ply > 0 && board.IsRepetition() || board.HalfMoveClock >= 100 {
		return 0
	}

	// check if the search should be stopped, time is checked every 10240 nodes
	nodes++

	var standPat int = eval.Eval(*cb)

	// found a better move
	if standPat > alpha {
		alpha = standPat
		// fails high
		if standPat >= beta {
			return beta
		}
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
