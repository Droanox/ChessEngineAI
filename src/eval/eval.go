package eval

import (
	"github.com/Droanox/ChessEngineAI/src/board"
)

// centipawn values corresponding to pieces
// pawn, knight, bishop, rook, queen, king in that order
// with the first row as white and second as black

// initValues initializes the piece values arrays for the mid game and end game
func initValues() {
	for piece := board.Pawn - 1; piece < board.King; piece++ {
		for square := 0; square < 64; square++ {
			tableMG[piece][square] = PieceValuesMG[piece] + piecesMG[piece][square^56]
			tableEG[piece][square] = PieceValuesEG[piece] + piecesEG[piece][square^56]
			tableMG[piece+6][square] = PieceValuesMG[piece] + piecesMG[piece][square]
			tableEG[piece+6][square] = PieceValuesEG[piece] + piecesEG[piece][square]
		}
	}
}

// initFileAndRanksMasks initializes the FileMasks and RankMasks arrays
func initMasks() {
	for square := 0; square < 64; square++ {
		FileMasks[square] = board.FileAOn << (square % 8)
		RankMasks[square] = board.Rank1On << ((square / 8) * 8)
		EastIsolatedMasks[square] = ((FileMasks[square] << 1) &^ board.FileAOn)
		WestIsolatedMasks[square] = ((FileMasks[square] >> 1) &^ board.FileHOn)
		KingFileMasks[square] = EastIsolatedMasks[square] | FileMasks[square] | WestIsolatedMasks[square]
		KingForwardSquares[board.White][square] = (board.KingAttacks[square] | board.KingAttacks[Min(square+8, 63)]) - board.KingAttacks[square]
		KingForwardSquares[board.Black][square] = (board.KingAttacks[square] | board.KingAttacks[Max(square-8, 0)]) - board.KingAttacks[square]
		PassedMasks[board.White][square] = ((EastIsolatedMasks[square] | WestIsolatedMasks[square]) ^ FileMasks[square]) << (8 + (8 * (square / 8)))
		PassedMasks[board.Black][square] = ((EastIsolatedMasks[square] | WestIsolatedMasks[square]) ^ FileMasks[square]) >> (8 + (8 * ((63 - square) / 8)))
	}
}

// Init initializes the evaluation package
func Init() {
	initValues()
	initMasks()
}

// Eval returns the evaluation of the given chess board
func Eval(cb board.ChessBoard) int {
	var pieceArr = []*uint64{
		0: &cb.WhitePawns, 1: &cb.WhiteKnights, 2: &cb.WhiteBishops, 3: &cb.WhiteRooks, 4: &cb.WhiteQueen, 5: &cb.WhiteKing,
		6: &cb.BlackPawns, 7: &cb.BlackKnights, 8: &cb.BlackBishops, 9: &cb.BlackRooks, 10: &cb.BlackQueen, 11: &cb.BlackKing,
	}

	var side, square, gamephase, attackingPiecesCount int
	var mg, eg, pawnsMissing [2]int
	var kingSquares, pawnAttackSquares, unsafeSquares [2]uint64
	var kingstart int

	// Calculate the number of pieces on the board
	var allPieces = cb.WhitePieces | cb.BlackPieces

	// calculate pawn attack squares
	pawnAttackSquares[0] = whitePawnAnyAttack(cb.WhitePawns)
	pawnAttackSquares[1] = blackPawnAnyAttack(cb.BlackPawns)

	// calcuate the number of pawns missing
	pawnsMissing[0] = 8 - board.BitCount(cb.WhitePawns)
	pawnsMissing[1] = 8 - board.BitCount(cb.BlackPawns)

	// calculate the king moveable squares
	kingstart = board.BitScanForward(cb.WhiteKing)
	kingSquares[0] = (board.KingAttacks[kingstart] &^ cb.WhitePieces) | KingForwardSquares[0][kingstart]
	kingstart = board.BitScanForward(cb.BlackKing)
	kingSquares[1] = (board.KingAttacks[kingstart] &^ cb.BlackPieces) | KingForwardSquares[1][kingstart]

	// calculate the unsafe squares for pieces
	unsafeSquares[0] = (cb.BlackPieces | pawnAttackSquares[board.Black])
	unsafeSquares[1] = (cb.WhitePieces | pawnAttackSquares[board.White])

	for i, pieceBoard := range pieceArr {
		side = pieceToColour[i]
		pieceBoardNew := pieceBoard
		oppositePieceBoard := pieceArr[((1-side)*6)+(i%6)]
		for bitboard := *pieceBoardNew; bitboard != board.EmptyBoard; bitboard &= bitboard - 1 {
			square = board.BitScanForward(bitboard)
			mg[side] += tableMG[i][square]
			eg[side] += tableEG[i][square]
			gamephase += gamephaseInc[i]

			switch i {
			case 0, 6: // Pawns
				// double pawn penalty
				if board.BitCount(FileMasks[square]&*pieceBoard) > 1 {
					mg[side] -= doublePawnPenaltyHalf[0]
					eg[side] -= doublePawnPenaltyHalf[1]
				}

				// isolated pawn penalty
				if EastIsolatedMasks[square]&*pieceBoard == 0 {
					mg[side] -= isolatedPawnPenaltyHalf[0]
					eg[side] -= isolatedPawnPenaltyHalf[1]
				}
				if WestIsolatedMasks[square]&*pieceBoard == 0 {
					mg[side] -= isolatedPawnPenaltyHalf[0]
					eg[side] -= isolatedPawnPenaltyHalf[1]
				}

				// Passed pawn bonus
				if PassedMasks[side][square]&*oppositePieceBoard == 0 && board.IndexMasks[square+offsetBySide[side]]&*pieceBoard == 0 {
					mg[side] += pastPawnBonus[(square^(56*side))/8]
					eg[side] += pastPawnBonus[(square^(56*side))/8]
				}
			case 1, 7: // Knights
				// decrease knight value as pawns decrease
				mg[side] -= pawnsMissing[1-side] * knightPawnPenalty
				eg[side] -= pawnsMissing[1-side] * knightPawnPenalty

				// calculate moveable squares
				knightAttacks := board.KnightAttacks[square] &^ unsafeSquares[side]

				// mobility bonus
				score := board.BitCount(knightAttacks) - 4
				mg[side] += score * knightMobility
				eg[side] += score * knightMobility

				// king attack update
				attackingPiecesCount += 2 * board.BitCount(knightAttacks&kingSquares[1-side])
			case 2, 8: // Bishops
				// bishop pair bonus
				if board.BitCount(*pieceBoard) == 2 {
					mg[side] += bishopPairBonus
					eg[side] += bishopPairBonus
				}

				// calculate moveable squares
				bishopAttacks := board.GetBishopAttacks(square, allPieces) &^ unsafeSquares[side]

				// mobility bonus
				score := board.BitCount(bishopAttacks) - 6
				mg[side] += score * bishopMobility
				eg[side] += score * bishopMobility

				// king attack update
				attackingPiecesCount += 2 * board.BitCount(bishopAttacks&kingSquares[1-side])
			case 3, 9: // Rooks
				// increase rook value as pawns decrease
				mg[side] += pawnsMissing[1-side] * rookPawnBonus
				eg[side] += pawnsMissing[1-side] * rookPawnBonus

				// rook on open file bonus
				if FileMasks[square]&(cb.WhitePawns|cb.BlackPawns) == 0 {
					mg[side] += rookOpenFile
					// give a bonus if there are multiple rooks or queens on the same file
					piecesOnFile := board.BitCount(FileMasks[square] & (*pieceBoard | *pieceArr[4+(side*6)]))
					if piecesOnFile > 1 {
						mg[side] += stackedPieceBonus * piecesOnFile
					}
				} else
				// rook on semi open file bonus
				if FileMasks[square]&*pieceArr[side*6] == 0 {
					mg[side] += rookSemiOpenFile
				}

				// rook on seventh rank bonus
				if (RankMasks[square])&rook7thRank[side] != 0 {
					pawnsOnSeventh := board.BitCount(rook7thRank[side] & *pieceArr[(1-side)*6])
					mg[side] += rookOnSeventh * pawnsOnSeventh
					eg[side] += rookOnSeventh * pawnsOnSeventh
				}

				// calcualte moveable squares
				rookAttacks := board.GetRookAttacks(square, allPieces) &^ unsafeSquares[side]

				// mobility bonus
				score := board.BitCount(rookAttacks) - 7
				mg[side] += score * rookMobility[0]
				eg[side] += score * rookMobility[1]

				// king attack update
				attackingPiecesCount += 3 * board.BitCount(rookAttacks&kingSquares[1-side])
			case 4, 10: // Queens
				queenAttacks := board.GetQueenAttacks(square, allPieces) &^ unsafeSquares[side]

				// king attack update
				attackingPiecesCount += 5 * board.BitCount(queenAttacks&kingSquares[1-side])
			case 5, 11: // Kings
				// king on open file penalty
				// if KingFileMasks[square]&(cb.WhitePawns|cb.BlackPawns) == 0 {
				// 	mg[side] -= kingOpenFile
				// 	if KingFileMasks[square]&(*pieceArr[3+((1-side)*6)]|*pieceArr[4+((1-side)*6)]) != 0 {
				// 		mg[side] -= kingOpenFile * 2
				// 	}
				// } else
				// // king on semi open file penalty
				// if KingFileMasks[square]&*pieceArr[side*6] == 0 {
				// 	mg[side] -= kingSemiOpenFile
				// }

				// attacking king zone
				// Idea from https://www.chessprogramming.org/King_Safety#Attacking_King_Zone
				mg[side] += SafetyTable[Min(attackingPiecesCount, 99)]
				// Add how many pieces are attacking the opponent for use in search.
				KingAttackingPieces[side] = attackingPiecesCount
				attackingPiecesCount = 0

				// pawn shield
				mg[side] += board.BitCount(board.KingAttacks[square]&*pieceArr[side*6]) * pawnMultiplier

				// king tropism
				mg[side] -= board.BitCount((board.GetQueenAttacks(square, allPieces) &^ unsafeSquares[side])) * kingTropismPenaltyMultiplier
			}
		}
	}

	// Caculate scores and phases
	var scoreMG int = mg[board.SideToMove] - mg[1-board.SideToMove]
	var scoreEG int = eg[board.SideToMove] - eg[1-board.SideToMove]
	var phaseMG int = gamephase
	if phaseMG > 24 {
		phaseMG = 24
	}
	var phaseEG int = 24 - phaseMG

	// Recalculate pawn value
	PawnValue = ((PieceValuesMG[0] * phaseMG) + (PieceValuesEG[0] * phaseEG)) / 24

	// adjust king attacking pieces depending on game phase
	KingAttackingPieces[0] = ((KingAttackingPieces[0] * phaseMG) + (1 * phaseEG)) / 24
	KingAttackingPieces[1] = ((KingAttackingPieces[1] * phaseMG) + (1 * phaseEG)) / 24

	// return tapered evaluation
	return (((scoreMG * phaseMG) + (scoreEG * phaseEG)) / 24)
}

// IsEndGame returns true if the given chess board is an end game
// isn't perfect, for quick computation
func IsEndGame(cb board.ChessBoard) bool {
	if board.SideToMove == board.White {
		return cb.WhitePieces&(cb.WhitePawns|cb.WhiteKing) == cb.WhitePieces
	} else {
		return cb.BlackPieces&(cb.BlackPawns|cb.BlackKing) == cb.BlackPieces
	}
}
