package engine

import (
	"github.com/Droanox/ChessEngineAI/src/board"
)

type TranspositionTable struct {
	key   uint64
	depth int
	score int
	flag  int
}

func ClearTT() {
	tt = [hashSize]TranspositionTable{}
}

// read the transposition table
// Implementation based on Bruce Moreland's implementation
// https://web.archive.org/web/20071031100051/http://www.brucemo.com/compchess/programming/hashing.htm
func ReadTT(alpha int, beta int, depth int) int {
	entryTT := &tt[board.HashKey%hashSize]

	if entryTT.key == board.HashKey {
		if entryTT.depth >= depth {
			var score int = entryTT.score

			// adjust the score based on the depth
			if score < MateScore {
				score += board.Ply
			}
			// adjust the score based on the depth
			if score > -MateScore {
				score -= board.Ply
			}

			if entryTT.flag == hashFlagExact {
				return score
			} else if entryTT.flag == hashFlagAlpha {
				if score <= alpha {
					return alpha
				}
			} else if entryTT.flag == hashFlagBeta {
				if score >= beta {
					return beta
				}
			}
		}
	}

	return noHash
}

func WriteTT(score int, depth int, flag int) {
	entryTT := &tt[board.HashKey%hashSize]

	// adjust the score based on the depth
	if score < MateScore {
		score -= board.Ply
	}
	// adjust the score based on the depth
	if score > -MateScore {
		score += board.Ply
	}

	entryTT.key = board.HashKey
	entryTT.depth = depth
	entryTT.score = score
	entryTT.flag = flag
}
