package eval

import "github.com/Droanox/ChessEngineAI/src/board"

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// whitePawnAnyAttack returns a bitboard of all squares attacked by the given white pawns
func whitePawnAnyAttack(pawns uint64) uint64 {
	return ((pawns << 9) & ^board.FileAOn) | ((pawns << 7) & ^board.FileHOn)
}

// blackPawnAnyAttack returns a bitboard of all squares attacked by the given black pawns
func blackPawnAnyAttack(pawns uint64) uint64 {
	return ((pawns >> 7) & ^board.FileAOn) | ((pawns >> 9) & ^board.FileHOn)
}

// // maskKingShield returns a bitboard of the squares attacked by a king
// func maskKingShieldWhite(square int) uint64 { // TODO
// 	var attacks, bitboard uint64 = 0, 0x100
// 	board.MoveBitsWhite(&bitboard, square)

// 	attacks = (bitboard<<1)&^board.FileAOn | (bitboard) | (bitboard>>1)&^board.FileHOn

// 	return attacks
// }

// // maskKingShield returns a bitboard of the squares attacked by a king
// func maskKingShieldBlack(square int) uint64 { // TODO
// 	var attacks, bitboard uint64 = 0, 0x80000000000000
// 	board.MoveBitsBlack(&bitboard, 63-square)

// 	attacks = (bitboard<<1)&^board.FileAOn | (bitboard) | (bitboard>>1)&^board.FileHOn

// 	return attacks
// }
