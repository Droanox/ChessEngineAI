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
