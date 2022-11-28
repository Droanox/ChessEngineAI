package board

func MaskPawnAttacks(square int, side int) uint64 { // TODO
	var attacks uint64
	var bitboard uint64
	SetBit(&bitboard, square)

	if side == White {
		if ((bitboard >> 7) & ^FileAOn) != uint64(0) {
			attacks |= bitboard >> 7
		}
		if ((bitboard >> 9) & ^FileHOn) != uint64(0) {
			attacks |= bitboard >> 9
		}
	}
	if side == Black {
		if ((bitboard << 9) & ^FileAOn) != uint64(0) {
			attacks |= bitboard << 9
		}
		if ((bitboard << 7) & ^FileHOn) != uint64(0) {
			attacks |= bitboard << 7
		}
	}

	return attacks
}
