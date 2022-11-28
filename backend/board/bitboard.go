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
// (a-8) represents the ranks
// (.) represents an empty square
// (X) represents a square with a piece on it
//

/* func logicalSum(bitmaps ...uint64) uint64 {
	var sum uint64
	for _, value := range bitmaps {
		sum = sum | value
	}
	return sum
} */

var (
	PawnAttacks   [2][64]uint64
	KnightAttacks [64]uint64
)

func AttackInit() {
	for i := 0; i < 64; i++ {
		PawnAttacks[0][i] = MaskPawnAttacks(i, 0)
		PawnAttacks[1][i] = MaskPawnAttacks(i, 1)
		KnightAttacks[i] = MaskKnightAttacks(i)
	}
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
