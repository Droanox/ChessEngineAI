package board

func MaskPawnAttacks(square int, side int) uint64 { // TODO
	var attacks uint64
	var bitboard uint64
	SetBit(&bitboard, square)

	if side == White {
		if ((bitboard >> 7) & ^FileAOn) != EmptyBoard {
			attacks |= bitboard >> 7
		}
		if ((bitboard >> 9) & ^FileHOn) != EmptyBoard {
			attacks |= bitboard >> 9
		}
	}
	if side == Black {
		if ((bitboard << 9) & ^FileAOn) != EmptyBoard {
			attacks |= bitboard << 9
		}
		if ((bitboard << 7) & ^FileHOn) != EmptyBoard {
			attacks |= bitboard << 7
		}
	}

	return attacks
}

func MaskKnightAttacks(square int) uint64 {
	var attacks uint64
	var bitboard uint64
	SetBit(&bitboard, square)

	if ((bitboard << 6) & (^FileGOn & ^FileHOn)) != EmptyBoard {
		attacks |= bitboard << 6
	}
	if ((bitboard << 15) & ^FileHOn) != EmptyBoard {
		attacks |= bitboard << 15
	}
	if ((bitboard << 17) & ^FileAOn) != EmptyBoard {
		attacks |= bitboard << 17
	}
	if ((bitboard << 10) & (^FileAOn & ^FileBOn)) != EmptyBoard {
		attacks |= bitboard << 10
	}

	if ((bitboard >> 6) & (^FileAOn & ^FileBOn)) != EmptyBoard {
		attacks |= bitboard >> 6
	}
	if ((bitboard >> 15) & ^FileAOn) != EmptyBoard {
		attacks |= bitboard >> 15
	}
	if ((bitboard >> 17) & ^FileHOn) != EmptyBoard {
		attacks |= bitboard >> 17
	}
	if ((bitboard >> 10) & (^FileGOn & ^FileHOn)) != EmptyBoard {
		attacks |= bitboard >> 10
	}
	return attacks
}

func MaskBishopAttacks(square int) uint64 {
	var attacks uint64
	nextRank := square / 8
	nextFile := square % 8

	for rank, file := nextRank+1, nextFile+1; rank < 8 && file < 8; rank, file = rank+1, file+1 {
		attacks |= uint64(1) << (rank*8 + file)
	}
	for rank, file := nextRank-1, nextFile+1; rank >= 0 && file < 8; rank, file = rank-1, file+1 {
		attacks |= uint64(1) << (rank*8 + file)
	}
	for rank, file := nextRank+1, nextFile-1; rank < 8 && file >= 0; rank, file = rank+1, file-1 {
		attacks |= uint64(1) << (rank*8 + file)
	}
	for rank, file := nextRank-1, nextFile-1; rank >= 0 && file >= 0; rank, file = rank-1, file-1 {
		attacks |= uint64(1) << (rank*8 + file)
	}

	return attacks
}

func MaskKingAttacks(square int) uint64 { // TODO
	var attacks uint64
	var bitboard uint64
	SetBit(&bitboard, square)

	if ((bitboard << 7) & ^FileHOn) != EmptyBoard {
		attacks |= bitboard << 7
	}
	if (bitboard << 8) != EmptyBoard {
		attacks |= bitboard << 8
	}
	if ((bitboard << 9) & ^FileAOn) != EmptyBoard {
		attacks |= bitboard << 9
	}
	if ((bitboard << 1) & ^FileAOn) != EmptyBoard {
		attacks |= bitboard << 1
	}

	if ((bitboard >> 7) & ^FileAOn) != EmptyBoard {
		attacks |= bitboard >> 7
	}
	if (bitboard >> 8) != EmptyBoard {
		attacks |= bitboard >> 8
	}
	if ((bitboard >> 9) & ^FileHOn) != EmptyBoard {
		attacks |= bitboard >> 9
	}
	if ((bitboard >> 1) & ^FileHOn) != EmptyBoard {
		attacks |= bitboard >> 1
	}
	return attacks
}
