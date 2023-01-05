package board

import "fmt"

//
// The bitboard format:
// 8|X|.|.|.|.|X|.|.|
// 7|.|.|.|.|.|.|.|.|
// 6|.|.|.|.|.|.|.|.|
// 5|.|.|.|X|.|.|X|.|
// 4|.|.|.|.|.|.|.|.|
// 3|.|.|.|.|.|.|.|.|
// 2|.|.|.|.|.|.|.|.|
// 1|.|.|.|.|X|.|.|.|
//   a b c d e f g h
//
// (a-h) represents the files
// (a-8) represents the ranks
// (.) represents an empty square
// (X) represents a square with a piece on it
//

type ChessBoard struct {
	WhitePawns   uint64
	WhiteKnights uint64
	WhiteBishops uint64
	WhiteRooks   uint64
	WhiteQueen   uint64
	WhiteKing    uint64

	BlackPawns   uint64
	BlackKnights uint64
	BlackBishops uint64
	BlackRooks   uint64
	BlackQueen   uint64
	BlackKing    uint64

	//whitePieces uint64
	//blackPieces uint64
}

var (
	pawnAttacks   [2][64]uint64
	knightAttacks [64]uint64
	bishopAttacks [64]uint64
	rookAttacks   [64]uint64
	queenAttacks  [64]uint64
	kingAttacks   [64]uint64
)

func attackInit() {
	for i := 0; i < 64; i++ {
		pawnAttacks[White][i] = maskPawnAttacks(i, White)
		pawnAttacks[Black][i] = maskPawnAttacks(i, Black)
		knightAttacks[i] = maskKnightAttacks(i)
		bishopAttacks[i] = maskBishopAttacks(i)
		rookAttacks[i] = maskRookAttacks(i)
		queenAttacks[i] = bishopAttacks[i] | rookAttacks[i]
		kingAttacks[i] = maskKingAttacks(i)
	}
}

func (cb *ChessBoard) Init() {
	attackInit()
	cb.parseFen(initialPositionFen)
	//whites := cb.WhiteRooks | cb.WhiteKnights | cb.WhiteBishops | cb.WhiteQueen | cb.WhiteKing | cb.WhitePawns
	//blacks := cb.BlackRooks | cb.BlackKnights | cb.BlackBishops | cb.BlackQueen | cb.BlackKing | cb.BlackPawns
	//PrintBitboard(whites | blacks)
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
