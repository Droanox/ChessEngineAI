package board

func whitePawnsEastAttacks(pawns uint64) uint64 {
	return (pawns << 9) & ^FileAOn
}
func whitePawnsWestAttacks(pawns uint64) uint64 {
	return (pawns << 7) & ^FileHOn
}
func blackPawnsEastAttacks(pawns uint64) uint64 {
	return (pawns >> 7) & ^FileAOn
}
func blackPawnsWestAttacks(pawns uint64) uint64 {
	return (pawns >> 9) & ^FileHOn
}

func whitePawnAnyAttack(pawns uint64) uint64 {
	return whitePawnsEastAttacks(pawns) | whitePawnsWestAttacks(pawns)
}
func blackPawnAnyAttack(pawns uint64) uint64 {
	return blackPawnsEastAttacks(pawns) | blackPawnsWestAttacks(pawns)
}

/*
func whitePawnsAbleToCaptureEast(wPawns uint64, bPawns uint64) uint64 {
	return wPawns & blackPawnsWestAttacks(bPawns)
}
func whitePawnsAbleToCaptureWest(wPawns uint64, bPawns uint64) uint64 {
	return wPawns & blackPawnsEastAttacks(bPawns)
}
func blackPawnsAbleToCaptureEast(wPawns uint64, bPawns uint64) uint64 {
	return bPawns & whitePawnsWestAttacks(wPawns)
}
func blackPawnsAbleToCaptureWest(wPawns uint64, bPawns uint64) uint64 {
	return bPawns & whitePawnsEastAttacks(wPawns)
}
*/

func maskPawnAttacks(side uint8, square int) uint64 { // TODO
	var bitboard uint64
	setBit(&bitboard, square)
	if side == White {
		bitboard = whitePawnAnyAttack(bitboard)
	} else {
		bitboard = blackPawnAnyAttack(bitboard)
	}

	return bitboard
}

func maskKnightAttacks(square int) uint64 {
	var attacks uint64
	var bitboard uint64
	setBit(&bitboard, square)

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

func maskKingAttacks(square int) uint64 { // TODO
	var attacks uint64
	var bitboard uint64
	setBit(&bitboard, square)

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
