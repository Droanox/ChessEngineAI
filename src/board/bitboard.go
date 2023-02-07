package board

import (
	"fmt"
)

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

	WhitePieces uint64
	BlackPieces uint64
}

var (
	pawnAttacks   [2][64]uint64
	knightAttacks [64]uint64
	bishopMasks   [64]uint64
	bishopAttacks [64][512]uint64
	rookMasks     [64]uint64
	rookAttacks   [64][4096]uint64
	kingAttacks   [64]uint64
	indexMasks    [64]uint64
)

func attackLeaperInit() {
	for i := 0; i < 64; i++ {
		pawnAttacks[White][i] = maskPawnAttacks(White, i)
		pawnAttacks[Black][i] = maskPawnAttacks(Black, i)
		knightAttacks[i] = maskKnightAttacks(i)
		kingAttacks[i] = maskKingAttacks(i)
		indexMasks[i] = 1 << i
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
	// Below init is used to inisialse the pawn, knight, and king pieces
	attackLeaperInit()

	// Below inits are used to initalise the bishop and rook attack tables
	attackSliderInit(true)
	attackSliderInit(false)
}

///////////////////////////////////////////////////////////////////
// Prints for debugging
///////////////////////////////////////////////////////////////////

func PrintBitboard(bitboard uint64) {
	for rank := 8; rank >= 1; rank-- {
		fmt.Print(rank)
		for file := 1; file <= 8; file++ {
			square := (rank-1)*8 + (file - 1)
			if ((bitboard >> square) & 1) != 0 {
				fmt.Print("|X")
			} else {
				fmt.Print("|·")
			}
		}
		fmt.Print("|\n")
	}
	fmt.Print("  a b c d e f g h\n")
}

func (cb ChessBoard) PrintChessBoard() {
	for rank := 8; rank >= 1; rank-- {
		fmt.Print(rank)
		for file := 1; file <= 8; file++ {
			square := (rank-1)*8 + (file - 1)
			fmt.Print(" ", pieceToASCII[cb.GetPieceInt(square)])
		}
		fmt.Print("\n")
	}
	fmt.Print("  a b c d e f g h\n")
	if SideToMove == White {
		fmt.Printf("\nSide to move: %s", "White")
	} else {
		fmt.Printf("\nSide to move: %s", "Black")
	}
	fmt.Printf("\nCastling rights (KQkq): ")
	if CastleRights&WhiteKingSide != 0 {
		fmt.Printf("K")
	} else {
		fmt.Printf("-")
	}
	if CastleRights&WhiteQueenSide != 0 {
		fmt.Printf("Q")
	} else {
		fmt.Printf("-")
	}
	if CastleRights&BlackKingSide != 0 {
		fmt.Printf("k")
	} else {
		fmt.Printf("-")
	}
	if CastleRights&BlackQueenSide != 0 {
		fmt.Printf("q")
	} else {
		fmt.Printf("-")
	}
	if Enpassant == -1 {
		fmt.Printf("\nEnpassant: %s", "-")
	} else {
		fmt.Printf("\nEnpassant: %s", IndexToSquare[Enpassant])
	}
	fmt.Printf("\nHalf move clock: %d", HalfMoveClock)
	fmt.Printf("\nFull move counter: %d\n\n", FullMoveCounter)
}

// Used to test whether the function IsAttackedBySide() returns the correct output
func (cb ChessBoard) PrintBitboardIsAttacked(side int) {
	for rank := 8; rank >= 1; rank-- {
		fmt.Print(rank)
		for file := 1; file <= 8; file++ {
			square := (rank-1)*8 + (file - 1)
			if cb.IsSquareAttackedBySide(square, side) {
				fmt.Print("|X")
			} else {
				fmt.Print("|.")
			}
		}
		fmt.Print("|\n")
	}
	fmt.Print("  a b c d e f g h\n")
}

// Used to print the hexadecimal bitboard representation,
// which is used to get bit masks
func PrintBitboardHex(bitboard uint64) {
	fmt.Printf("%s\n", fmt.Sprintf("0x%X", bitboard))
}

func PrintMove(move Move) {
	fmt.Printf("%4s%s%-6s%-10s%-12s%04b%9d\n", "",
		IndexToSquare[move.GetMoveStart()],
		IndexToSquare[move.GetMoveEnd()],
		IntToPiece[move.GetMoveStartPiece()],
		IntToPiece[move.GetMoveEndPiece()],
		move.GetMoveFlags(),
		move.Index)
}
func PrintMoveList(move []Move) {
	fmt.Printf("\n%4s%-8s%-10s%-12s%04s%8s\n\n",
		"", "Move", "Piece", "Captured", "Flags", "Index")
	length := len(move)
	for i := 0; i < length; i++ {
		PrintMove(move[i])
	}
	fmt.Printf("\n")
}
