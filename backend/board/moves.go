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
