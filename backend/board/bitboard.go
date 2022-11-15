package board

import "fmt"

//
// The bitboard format:
// 8|.|.|.|.|.|.|.|.|
// 7|.|.|.|.|.|.|.|.|
// 6|.|.|.|.|.|.|.|.|
// 5|.|.|.|.|.|.|.|.|
// 4|.|.|.|.|.|.|.|.|
// 3|.|.|.|.|.|.|.|.|
// 2|.|.|.|.|.|.|.|.|
// 1|.|.|.|.|.|.|.|.|
//   a b c d e f g h
//
// (a-h) represents the files
// (a-8) represents the rank
// (.) represents an empty square
// (X) represents a square with a piece on
//

/* func logicalSum(bitmaps ...uint64) uint64 {
	var sum uint64
	for _, value := range bitmaps {
		sum = sum | value
	}
	return sum
} */

var (
//pawnAttacks [2][64]uint64
)

func MaskPawnAttacks(square int, side int) uint64 { // TODO
	var attack uint64
	var bitboard uint64

	SetBit(&bitboard, square)

	return attack
}

func PrintBitboard(bitboard uint64) {
	for rank := 8; rank >= 1; rank-- {
		fmt.Print(rank)
		for file := 1; file <= 8; file++ {
			square := (rank-1)*8 + (file - 1)
			if ((bitboard >> square) & 1) != 0 {
				fmt.Print("|X")
			} else {
				fmt.Print("|.")
			}
		}
		fmt.Print("|\n")
	}
	fmt.Print("  a b c d e f g h\n")
}
