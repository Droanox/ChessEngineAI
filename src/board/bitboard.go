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

// ChessBoard is the main data structure for the chess engine
// Containing the bitboards for each piece and collections of pieces from both sides
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

// The following are used to initialise the attack tables
var (
	PawnAttacks   [2][64]uint64
	KnightAttacks [64]uint64
	bishopMasks   [64]uint64
	bishopAttacks [64][512]uint64
	rookMasks     [64]uint64
	rookAttacks   [64][4096]uint64
	KingAttacks   [64]uint64
	IndexMasks    [64]uint64
)

// The following are used to initialise the attack tables
func attackLeaperInit() {
	for i := 0; i < 64; i++ {
		PawnAttacks[White][i] = maskPawnAttacks(White, i)
		PawnAttacks[Black][i] = maskPawnAttacks(Black, i)
		KnightAttacks[i] = maskKnightAttacks(i)
		KingAttacks[i] = maskKingAttacks(i)
		IndexMasks[i] = 1 << i
	}
}

// The following are used to initialise the attack tables
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

// The following are used to initialise the attack tables
func (cb *ChessBoard) Init() {
	// Below init is used to inisialse the pawn, knight, and king pieces
	attackLeaperInit()

	// Below inits are used to initalise the bishop and rook attack tables
	attackSliderInit(true)
	attackSliderInit(false)

	// Below init is used to initialise the hash keys for zobrist hashing
	initHash()
}

///////////////////////////////////////////////////////////////////
// Prints for debugging
///////////////////////////////////////////////////////////////////

// PrintBitboard prints the bitboard in a human readable format
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

// PrintChessBoard prints the chessboard in a human readable format
// Containing pieces with the ASCII representation
// and all the other ChessBoard Info
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
	if CastleRights&whiteKingSide != 0 {
		fmt.Printf("K")
	} else {
		fmt.Printf("-")
	}
	if CastleRights&whiteQueenSide != 0 {
		fmt.Printf("Q")
	} else {
		fmt.Printf("-")
	}
	if CastleRights&blackKingSide != 0 {
		fmt.Printf("k")
	} else {
		fmt.Printf("-")
	}
	if CastleRights&blackQueenSide != 0 {
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

// used by PrintMoveList() to print consecutive moves
func printMove(move Move) {
	fmt.Printf("%4s%s%-6s%-10s%-12s%04b%5s%-4d", "",
		IndexToSquare[move.GetMoveStart()],
		IndexToSquare[move.GetMoveEnd()],
		IntToPiece[move.GetMoveStartPiece()],
		IntToPiece[move.GetMoveCapturedPiece()],
		move.GetMoveFlags(),
		"",
		move.Score)
}

// Used to print the move list in a human readable format
// along with the index of the move in the move list
func PrintMoveList(move []Move) {
	fmt.Printf("\n%4s%-8s%-10s%-12s%04s%9s%8s\n\n",
		"", "Move", "Piece", "Captured", "Flags", "Score", "Index")
	length := len(move)
	for i := 0; i < length; i++ {
		printMove(move[i])
		fmt.Printf("%9d\n", i)
	}
	fmt.Printf("\n")
}
