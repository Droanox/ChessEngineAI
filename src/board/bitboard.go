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
}

var (
	pawnAttacks   [2][64]uint64
	knightAttacks [64]uint64
	bishopMasks   [64]uint64
	bishopAttacks [64][512]uint64
	rookMasks     [64]uint64
	rookAttacks   [64][4096]uint64
	kingAttacks   [64]uint64
)

func attackLeaperInit() {
	for i := 0; i < 64; i++ {
		pawnAttacks[White][i] = maskPawnAttacks(White, i)
		pawnAttacks[Black][i] = maskPawnAttacks(Black, i)
		knightAttacks[i] = maskKnightAttacks(i)
		kingAttacks[i] = maskKingAttacks(i)
	}
}

func attackSliderInit(isBishop bool) {
	for square := 0; square < 64; square++ {
		bishopMasks[square] = maskMagicBishopAttacks(square)
		rookMasks[square] = maskMagicRookAttacks(square)
		var mask uint64
		if isBishop {
			mask = bishopMasks[square]
		} else {
			mask = rookMasks[square]
		}
		var bitsCounted int = BitCount(mask)
		var occupancyIndex int = 1 << bitsCounted

		for j := 0; j < occupancyIndex; j++ {
			var occupancy uint64 = setMagicOccupancies(j, bitsCounted, mask)
			if isBishop {
				var magicIndex int = int((occupancy * bishopMagicNumber[square]) >> (64 - bishopBits[square]))
				bishopAttacks[square][magicIndex] = maskBishopAttacks(square, occupancy)
			} else {
				var magicIndex int = int((occupancy * rookMagicNumber[square]) >> (64 - rookBits[square]))
				rookAttacks[square][magicIndex] = maskRookAttacks(square, occupancy)
			}
		}
	}
}

func (cb *ChessBoard) Init() {
	attackLeaperInit()
	attackSliderInit(true)
	attackSliderInit(false)
	//Below init is used to get the first iteration of magic number,
	//It is a set that works, it may not be the best set
	//MagicInit()
	cb.parseFen(initialPositionFen)
}

func (cb *ChessBoard) Test() {
	var block uint64
	//whites := cb.WhiteRooks | cb.WhiteKnights | cb.WhiteBishops | cb.WhiteQueen | cb.WhiteKing | cb.WhitePawns
	//blacks := cb.BlackRooks | cb.BlackKnights | cb.BlackBishops | cb.BlackQueen | cb.BlackKing | cb.BlackPawns
	setBit(&block, 0)
	setBit(&block, 9)
	setBit(&block, 18)
	setBit(&block, 36)
	PrintBitboard(GetBishopAttacks(27, block))
	//PrintBitboard((maskMagicRookAttacks(23) * rookMagicNumber[23]) >> (64 - rookBits[23]))
	//PrintBitboard(whites | blacks)
	//PrintBitboard(maskMagicRookAttacks(0))
	//PrintBitboard(generateMagicNumberCandidate())
	//PrintBitboard(generateMagicNumberCandidate())
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
