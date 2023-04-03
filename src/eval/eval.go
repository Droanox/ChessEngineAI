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
			tableMG[piece][square] = pieceValuesMG[piece] + piecesMG[piece][square^56]
			tableEG[piece][square] = pieceValuesEG[piece] + piecesEG[piece][square^56]
			tableMG[piece+6][square] = pieceValuesMG[piece] + piecesMG[piece][square]
			tableEG[piece+6][square] = pieceValuesEG[piece] + piecesEG[piece][square]
		}
	}
}

// initFileAndRanksMasks initializes the FileMasks and RankMasks arrays
func initMasks() {
	for square := 0; square < 64; square++ {
		FileMasks[square] = board.FileAOn << (square % 8)
		RankMasks[square] = board.Rank1On << ((square / 8) * 8)
		IsolatedMasks[square] = ((FileMasks[square] << 1) &^ board.FileAOn) | ((FileMasks[square] >> 1) &^ board.FileHOn)
		PassedMasks[board.White][square] = (IsolatedMasks[square] ^ FileMasks[square]) << (8 + (8 * (square / 8)))
		PassedMasks[board.Black][square] = (IsolatedMasks[square] ^ FileMasks[square]) >> (8 + (8 * ((63 - square) / 8)))
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

	var side, square, gamephase int
	var mg, eg [2]int

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
					mg[side] += doublePawnPenalty
					eg[side] += doublePawnPenalty
					// fmt.Println("Double count penalty on:", "doublePawns", board.IndexToSquare[square], doublePawnPenalty)
				}
				// isolated pawn penalty
				if IsolatedMasks[square]&*pieceBoard == 0 {
					mg[side] += isolatedPawnPenalty
					eg[side] += isolatedPawnPenalty
					// fmt.Println("Isolated pawn penalty on:", board.IndexToSquare[square], isolatedPawnPenalty)
				}
				// Passed pawn bonus
				if PassedMasks[side][square]&*pieceArr[(1-side)*6] == 0 {
					mg[side] += PastPawnBonus[(square^(56*side))/8]
					eg[side] += PastPawnBonus[(square^(56*side))/8]
					// fmt.Println("Passed pawn bonus on:", board.IndexToSquare[square], PastPawnBonus[(square^(56*side))/8])
				}
			case 1, 7: // Knights
				// mobility bonus
				score := board.BitCount(board.KnightAttacks[square]&^*allPieceArr[side]) - 4
				mg[side] += score * knightMobility
				eg[side] += score * knightMobility
			case 2, 8: // Bishops
				// mobility bonus
				score := board.BitCount(board.GetBishopAttacks(square, cb.WhitePieces|cb.BlackPieces)) - 6
				mg[side] += score * bishopMobility
				eg[side] += score * bishopMobility
			case 3, 9: // Rooks
				// rook on open file bonus
				if FileMasks[square]&(cb.WhitePawns|cb.BlackPawns) == 0 {
					mg[side] += openFile[0]
					eg[side] += openFile[0]
					// fmt.Println("Rook on open file bonus on:", board.IndexToSquare[square], openFile)
				}
				// rook on semi open file bonus
				if FileMasks[square]&*pieceArr[side*6] == 0 {
					mg[side] += semiOpenFile[0]
					eg[side] += semiOpenFile[0]
					// fmt.Println("Rook on semi open file bonus on:", board.IndexToSquare[square], semiOpenFile)
				}
				// mobility bonus
				score := board.BitCount(board.GetRookAttacks(square, cb.WhitePieces|cb.BlackPieces)) - 7
				mg[side] += score * rookMobility[0]
				eg[side] += score * rookMobility[1]

			case 4, 10: // Queens
				score := board.BitCount(board.GetQueenAttacks(square, cb.WhitePieces|cb.BlackPieces))
				mg[side] += score
				eg[side] += score
			case 5, 11: // Kings
				// king on open file penalty
				if FileMasks[square]&(cb.WhitePawns|cb.BlackPawns) == 0 {
					mg[side] -= openFile[1]
					eg[side] -= openFile[2]
					// fmt.Println("King on open file penalty on:", board.IndexToSquare[square], openFile)
				}
				// king on semi open file penalty
				if FileMasks[square]&*pieceArr[side*6] == 0 {
					mg[side] -= semiOpenFile[1]
					eg[side] -= semiOpenFile[2]
					// fmt.Println("King on semi open file penalty on:", board.IndexToSquare[square], semiOpenFile)
				}
				// king safety bonus
				score := board.BitCount(board.KingAttacks[square]&*allPieceArr[side]) * kingSafetyBonus
				mg[side] += score
				eg[side] += score
			}
		}
	}

	var scoreMG int = mg[board.SideToMove] - mg[1-board.SideToMove]
	var scoreEG int = eg[board.SideToMove] - eg[1-board.SideToMove]
	var phaseMG int = gamephase
	if phaseMG > 24 {
		phaseMG = 24
	}
	var phaseEG int = 24 - phaseMG
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
