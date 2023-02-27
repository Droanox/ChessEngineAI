package board

// whitePawnAnyAttack returns a bitboard of all squares attacked by the given white pawns
func whitePawnAnyAttack(pawns uint64) uint64 {
	return ((pawns << 9) & ^FileAOn) | ((pawns << 7) & ^FileHOn)
}

// blackPawnAnyAttack returns a bitboard of all squares attacked by the given black pawns
func blackPawnAnyAttack(pawns uint64) uint64 {
	return ((pawns >> 7) & ^FileAOn) | ((pawns >> 9) & ^FileHOn)
}

// maskPawnAttacks returns a bitboard of the square attacks by a given side's pawn
// idea from https://www.chessprogramming.org/Pawn_Attacks_(Bitboards)
func maskPawnAttacks(side uint8, square int) uint64 { // TODO
	var bitboard uint64
	setBit(&bitboard, square)
	if side == White {
		return whitePawnAnyAttack(bitboard)
	} else {
		return blackPawnAnyAttack(bitboard)
	}
}

// maskKnightAttacks returns a bitboard of the squares attacked by a knight
// idea from https://www.chessprogramming.org/Knight_Pattern
func maskKnightAttacks(square int) uint64 {
	var bitboard, half1, half2 uint64
	setBit(&bitboard, square)

	half1 = ((bitboard >> 1) & ^FileHOn) | ((bitboard << 1) & ^FileAOn)
	half2 = ((bitboard >> 2) & ^FileGHOn) | ((bitboard << 2) & ^FileABOn)

	return (half1 << 16) | (half1 >> 16) | (half2 << 8) | (half2 >> 8)
}

// maskKingAttacks returns a bitboard of the squares attacked by a king
// idea from https://www.chessprogramming.org/King_Pattern
func maskKingAttacks(square int) uint64 { // TODO
	var attacks, bitboard uint64
	setBit(&bitboard, square)

	attacks = ((bitboard << 1) &^ FileAOn) | ((bitboard >> 1) &^ FileHOn)
	bitboard |= attacks
	attacks |= (bitboard << 8) | (bitboard >> 8)

	return attacks
}
