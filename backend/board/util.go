package board

func SetBit(bitboard *uint64, square int) {
	*bitboard |= (1 << uint64(square))
}
