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
		KingSquares[board.White][square] = board.KingAttacks[square] // | board.KingAttacks[Min(square+8, 63)]
		KingSquares[board.Black][square] = board.KingAttacks[square] // | board.KingAttacks[Max(square-8, 0)]
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
	var allPieceArr = []*uint64{
		0: &cb.WhitePieces, 1: &cb.BlackPieces,
	}

	var side, square, gamephase, attackingPiecesCount, valueOfAttacks int
	var mg, eg, pawnsMissing [2]int
	var kingSquares, pawnAttackSquares, unsafeSquares [2]uint64

	// Calculate the number of pieces on the board
	var allPieces = cb.WhitePieces | cb.BlackPieces

	// calculate pawn attack squares
	pawnAttackSquares[0] = whitePawnAnyAttack(cb.WhitePawns)
	pawnAttackSquares[1] = blackPawnAnyAttack(cb.BlackPawns)

	// calcuate the number of pawns missing
	pawnsMissing[0] = 8 - board.BitCount(cb.WhitePawns)
	pawnsMissing[1] = 8 - board.BitCount(cb.BlackPawns)

	// calculate the king moveable squares
	kingSquares[0] = KingSquares[0][board.BitScanForward(cb.WhiteKing)] &^ allPieces
	kingSquares[1] = KingSquares[1][board.BitScanForward(cb.BlackKing)] &^ allPieces

	// calculate the safe squares for pieces
	unsafeSquares[0] = (*allPieceArr[board.Black] | pawnAttackSquares[board.Black])
	unsafeSquares[1] = (*allPieceArr[board.White] | pawnAttackSquares[board.White])

	for i, pieceBoard := range pieceArr {
		pieceBoardNew := pieceBoard
		side = pieceToColour[i]
		for bitboard := *pieceBoardNew; bitboard != board.EmptyBoard; bitboard &= bitboard - 1 {
			square = board.BitScanForward(bitboard)
			mg[side] += tableMG[i][square]
			eg[side] += tableEG[i][square]
			gamephase += gamephaseInc[i]

			switch i {
			case 0, 6: // Pawns
				// double pawn penalty
				if board.BitCount(FileMasks[square]&*pieceBoard) > 1 {
					mg[side] -= doublePawnPenaltyHalf
					eg[side] -= doublePawnPenaltyHalf
					// fmt.Println("Double count penalty on:", "doublePawns", board.IndexToSquare[square], doublePawnPenalty)
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
				if PassedMasks[side][square]&*pieceArr[(1-side)*6] == 0 && *pieceArr[(side)*6]&board.IndexMasks[square+offsetBySide[side]] == 0 {
					mg[side] += pastPawnBonus[(square^(56*side))/8] / 2
					eg[side] += pastPawnBonus[(square^(56*side))/8]
					// fmt.Println("Passed pawn bonus on:", board.IndexToSquare[square], PastPawnBonus[(square^(56*side))/8])
				}
			case 1, 7: // Knights
				// decrease knight value as pawns decrease
				mg[side] -= pawnsMissing[1-side] * knightPawnPenalty
				eg[side] -= pawnsMissing[1-side] * knightPawnPenalty

				// calcualte moveable squares
				knightAttacks := board.KnightAttacks[square] &^ unsafeSquares[side]

				// mobility bonus
				score := board.BitCount(knightAttacks) - 4
				mg[side] += score * knightMobility
				eg[side] += score * knightMobility

				// king attack update
				numAttacks := board.BitCount(knightAttacks&kingSquares[1-side]) * attackerWeights[1]
				if numAttacks > 0 {
					valueOfAttacks += numAttacks
					attackingPiecesCount++
				}
			case 2, 8: // Bishops
				// bishop pair bonus
				if board.BitCount(*pieceArr[2+(side*6)]) == 2 {
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
				numAttacks := board.BitCount(bishopAttacks&kingSquares[1-side]) * attackerWeights[2]
				if numAttacks > 0 {
					valueOfAttacks += numAttacks
					attackingPiecesCount++
				}
			case 3, 9: // Rooks
				// increase rook value as pawns decrease
				mg[side] += pawnsMissing[1-side] * rookPawnBonus
				eg[side] += pawnsMissing[1-side] * rookPawnBonus

				// rook on open file bonus
				if FileMasks[square]&(cb.WhitePawns|cb.BlackPawns) == 0 {
					mg[side] += openFile
					eg[side] += openFile / 2
					// fmt.Println("Rook on open file bonus on:", board.IndexToSquare[square], openFile)
					piecesOnFile := board.BitCount(FileMasks[square] & (*pieceArr[3+(side*6)] | *pieceArr[4+(side*6)]))
					if piecesOnFile > 1 {
						mg[side] += stackedPieceBonus * piecesOnFile
					}
				} else
				// rook on semi open file bonus
				if FileMasks[square]&*pieceArr[side*6] == 0 {
					mg[side] += semiOpenFile
					eg[side] += semiOpenFile / 2
					// fmt.Println("Rook on semi open file bonus on:", board.IndexToSquare[square], semiOpenFile)
				}

				// rook on seventh rank bonus
				if (RankMasks[square])&rook7thRank[side] != 0 {
					pawnsOnSeventh := board.BitCount(board.GetRookAttacks(square, allPieces) & *pieceArr[(1-side)*6])
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
				numAttacks := board.BitCount(rookAttacks&kingSquares[1-side]) * attackerWeights[3]
				if numAttacks > 0 {
					valueOfAttacks += numAttacks
					attackingPiecesCount++
				}
			case 4, 10: // Queens
				queenAttacks := board.GetQueenAttacks(square, allPieces) &^ unsafeSquares[side]

				// king attack update
				numAttacks := board.BitCount(queenAttacks&kingSquares[1-side]) * attackerWeights[4]
				if numAttacks > 0 {
					valueOfAttacks += numAttacks
					attackingPiecesCount++
				}
			case 5, 11: // Kings
				// king on open file penalty
				if KingFileMasks[square]&(cb.WhitePawns|cb.BlackPawns) == 0 {
					mg[side] -= openFile * 2
					// fmt.Println("King on open file penalty on:", board.IndexToSquare[square], openFile)
				} else
				// king on semi open file penalty
				if KingFileMasks[square]&*pieceArr[side*6] == 0 {
					mg[side] -= semiOpenFile * 2
					// fmt.Println("King on semi open file penalty on:", board.IndexToSquare[square], semiOpenFile)
				}

				// attacking king zone
				// Idea from https://www.chessprogramming.org/King_Safety#Attacking_King_Zone
				mg[side] += (valueOfAttacks * numberOfAttacksWeight[Min(attackingPiecesCount, 6)]) / 100

				// pawn shield
				mg[side] += board.BitCount(board.KingAttacks[square]&*pieceArr[side*6]) * pawnMultiplier[0]

				// king tropism
				mg[side] -= board.BitCount((board.GetQueenAttacks(square, allPieces) &^ unsafeSquares[side])) * kingTropismPenaltyMultiplier

				attackingPiecesCount, valueOfAttacks = 0, 0
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

	// return tapered evaluation
	return ((scoreMG * phaseMG) + (scoreEG * phaseEG)) / 24
}

// IsEndGame returns true if the given chess board is an end game
// isn't perfect, for quick computation
func IsEndGame(cb board.ChessBoard) bool {
	if board.SideToMove == board.White {
		return (cb.WhiteRooks|cb.WhiteQueen) == 0 && cb.WhitePawns != 0
	} else {
		return (cb.BlackRooks|cb.BlackQueen) == 0 && cb.BlackPawns != 0
	}
}
