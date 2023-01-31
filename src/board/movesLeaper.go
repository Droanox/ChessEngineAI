package board

func whitePawnAnyAttack(pawns uint64) uint64 {
	return ((pawns << 9) & ^FileAOn) | ((pawns << 7) & ^FileHOn)
}
func blackPawnAnyAttack(pawns uint64) uint64 {
	return ((pawns >> 7) & ^FileAOn) | ((pawns >> 9) & ^FileHOn)
}

func maskPawnAttacks(side uint8, square int) uint64 { // TODO
	var bitboard uint64
	setBit(&bitboard, square)
	if side == White {
		return whitePawnAnyAttack(bitboard)
	} else {
		return blackPawnAnyAttack(bitboard)
	}
}

func maskKnightAttacks(square int) uint64 {
	var bitboard, half1, half2 uint64
	setBit(&bitboard, square)

	half1 = ((bitboard >> 1) & ^FileHOn) | ((bitboard << 1) & ^FileAOn)
	half2 = ((bitboard >> 2) & ^FileGHOn) | ((bitboard << 2) & ^FileABOn)

	return (half1 << 16) | (half1 >> 16) | (half2 << 8) | (half2 >> 8)
}

func maskKingAttacks(square int) uint64 { // TODO
	var attacks, bitboard uint64
	setBit(&bitboard, square)

	attacks = ((bitboard << 1) &^ FileAOn) | ((bitboard >> 1) &^ FileHOn)
	bitboard |= attacks
	attacks |= (bitboard << 8) | (bitboard >> 8)

	return attacks
}
