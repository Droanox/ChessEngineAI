package engine

import (
	"github.com/Droanox/ChessEngineAI/src/board"
)

// move ordering
// https://www.chessprogramming.org/Move_Ordering
// Give each move a score depending on the type of move
func scoreMoves(movelist *[]board.Move, bestMove board.Move) {
	moves := (*movelist)

	for i := range moves {
		flags := moves[i].GetMoveFlags()
		switch {
		case moves[i].Move == bestMove.Move:
			moves[i].Score = moveOrderOffset + 72
		case flags == board.MoveQueenPromotionCapture:
			moves[i].Score = moveOrderOffset + 71
		case flags == board.MoveQueenPromotion:
			moves[i].Score = moveOrderOffset + 70
		case flags&board.MoveCaptures != 0:
			moves[i].Score = moveOrderOffset + MVV_LVA[moves[i].GetMoveCapturedPiece()][moves[i].GetMoveStartPiece()]
		case moves[i].Move == killerMoves[0][board.Ply].Move:
			moves[i].Score = moveOrderOffset - 1
		case board.Ply > 1 && moves[i].Move == killerMoves[0][board.Ply-2].Move:
			moves[i].Score = moveOrderOffset - 2
		case moves[i].Move == killerMoves[1][board.Ply].Move:
			moves[i].Score = moveOrderOffset - 3
		case board.Ply > 1 && moves[i].Move == killerMoves[1][board.Ply-2].Move:
			moves[i].Score = moveOrderOffset - 4
		case flags == board.MoveKingCastle:
			moves[i].Score = moveOrderOffset - 5
		case flags == board.MoveQueenCastle:
			moves[i].Score = moveOrderOffset - 6
		case flags >= board.MoveKnightPromotion:
			moves[i].Score = moveOrderOffset - 7
		default:
			if bfScore[moves[i].GetMoveStart()][moves[i].GetMoveEnd()] > 0 {
				moves[i].Score = hhScore[moves[i].GetMoveStart()][moves[i].GetMoveEnd()] / bfScore[moves[i].GetMoveStart()][moves[i].GetMoveEnd()]
			}

			// if board.Ply > 0 && moves[i].Move == counterMoves[movesMade[board.Ply-1].GetMoveStart()][movesMade[board.Ply-1].GetMoveEnd()] {
			// 	moves[i].Score += 5
			// }
		}
	}
}

func pickMove(movelist *[]board.Move, currentIndex int) {
	var moves = (*movelist)

	// Move through the list and if there's a move greater than the current Index
	// swap the move, keep doing this until a cut off occurs
	for nextMove := currentIndex + 1; nextMove < len(moves); nextMove++ {
		if (moves)[nextMove].Score > (moves)[currentIndex].Score {
			// swap elements in movelist
			(moves)[nextMove], (moves)[currentIndex] = (moves)[currentIndex], (moves)[nextMove]
		}
	}
}

/*
package engine

import (
	"github.com/Droanox/ChessEngineAI/src/board"
)

func scoreMoves(movelist *[]board.Move, bestMove board.Move) {
	moves := (*movelist)


		if pvFollowed {
			pvFollowed = false
			for i := 0; i < len(moves); i++ {
				if moves[i].Move == pvTable[0][board.Ply].Move {
					moves[i].Score = moveOrderOffset + 150
					pvFollowed = true
					return
				}
			}
		}

	for i := 0; i < len(moves); i++ {
		flags := moves[i].GetMoveFlags()

		if bestMove.Move == moves[i].Move {
			moves[i].Score = moveOrderOffset + 160
		} else if flags&board.MoveCaptures != 0 {
			moves[i].Score = moveOrderOffset + MVV_LVA[moves[i].GetMoveCapturedPiece()][moves[i].GetMoveStartPiece()]
		} else {
			switch {
			case moves[i].Move == mateKillerMoves[board.Ply].Move:
				// mate killer move
				moves[i].Score = moveOrderOffset - 10
			case moves[i].Move == killerMoves[0][board.Ply].Move:
				// first killer move
				moves[i].Score = moveOrderOffset - 20
			case moves[i].Move == killerMoves[1][board.Ply].Move:
				// second killer move
				moves[i].Score = moveOrderOffset - 30
			case flags&board.MoveKingCastle != 0:
				// king castle
				moves[i].Score = moveOrderOffset - 40
			case flags&board.MoveQueenCastle != 0:
				// queen castle
				moves[i].Score = moveOrderOffset - 50
			default:
				// history move
				moves[i].Score = historyMoves[(moves[i].GetMoveStartPiece()+(board.SideToMove*6))-1][moves[i].GetMoveEnd()]
			}
			// if i > 0 && counterMoves[board.SideToMove][moves[i-1].GetMoveStart()][moves[i-1].GetMoveEnd()].Move == moves[i].Move {
			// moves[i].Score += 1
			// }
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
*/
