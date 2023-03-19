package board

import (
	"math/rand"
)

// TranspositionTable is a hash table that stores the best move and score for a given position
// It is used to avoid searching the same position multiple times
func initHash() {
	rand.Seed(9052375092375182470)

	// loop through all the pieces on the board
	for piece := 0; piece < 12; piece++ {
		// loop through all the squares on the board
		for square := 0; square < 64; square++ {
			// generate a random number for each piece on each square
			pieceKeys[piece][square] = rand.Uint64()
		}
	}
	for square := 0; square < 64; square++ {
		// generate a random number for each enpassant square
		enpassantKeys[square] = rand.Uint64()
	}
	for castleSide := 0; castleSide < 16; castleSide++ {
		// generate a random number for each castle side
		castleKeys[castleSide] = rand.Uint64()
	}

	// generate a random number for a side
	sideKey = rand.Uint64()
}

func GenHash(cb ChessBoard) uint64 {
	finalHash := uint64(0)

	pieceArr := []uint64{
		1: cb.WhitePawns, 2: cb.WhiteKnights, 3: cb.WhiteBishops, 4: cb.WhiteRooks, 5: cb.WhiteQueen, 6: cb.WhiteKing,
		7: cb.BlackPawns, 8: cb.BlackKnights, 9: cb.BlackBishops, 10: cb.BlackRooks, 11: cb.BlackQueen, 12: cb.BlackKing,
	}

	for piece := WhitePawns; piece <= BlackKing; piece++ {
		for bitboard := pieceArr[piece]; bitboard != EmptyBoard; bitboard &= bitboard - 1 {
			square := BitScanForward(bitboard)
			finalHash ^= pieceKeys[piece-1][square]
		}
	}
	if Enpassant != 64 {
		finalHash ^= enpassantKeys[Enpassant]
	}
	finalHash ^= castleKeys[CastleRights]
	if SideToMove == Black {
		finalHash ^= sideKey
	}

	return finalHash
}
