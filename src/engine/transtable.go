package engine

import (
	"math/rand"

	"github.com/Droanox/ChessEngineAI/src/board"
)

// TranspositionTable is a hash table that stores the best move and score for a given position
// It is used to avoid searching the same position multiple times
func initHash() {
	// loop through all the pieces on the board
	for piece := board.WhitePawns - 1; piece < board.BlackKing; piece++ {
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

func GenHash(cb board.ChessBoard) uint64 {
	finalHash := uint64(0)

	pieceArr := []uint64{
		1: cb.WhitePawns, 2: cb.WhiteKnights, 3: cb.WhiteBishops, 4: cb.WhiteRooks, 5: cb.WhiteQueen, 6: cb.WhiteKing,
		7: cb.BlackPawns, 8: cb.BlackKnights, 9: cb.BlackBishops, 10: cb.BlackRooks, 11: cb.BlackQueen, 12: cb.BlackKing,
	}

	for piece := board.WhitePawns; piece < board.BlackKing; piece++ {
		for bitboard := pieceArr[piece]; bitboard != board.EmptyBoard; bitboard &= bitboard - 1 {
			square := board.BitScanForward(bitboard)
			finalHash ^= pieceKeys[piece][square]
		}
	}
	if board.Enpassant != 64 {
		finalHash ^= enpassantKeys[board.Enpassant]
	}
	finalHash ^= castleKeys[board.CastleRights]
	if board.SideToMove == board.Black {
		finalHash ^= sideKey
	}

	return finalHash
}
