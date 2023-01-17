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

	// Below init is used to initalise the bishop and rook attack tables
	attackSliderInit(true)
	attackSliderInit(false)

	// Below init is used to get the first iteration of magic number, It is
	// a set that works, it may not be the best set, only need to do it once,
	// and output is used as a variable
	// MagicInit()

	cb.parseFen("1pbr2k1/P5p1/4p2p/5pb1/N7/1R6/3p1PPP/R1P3K1 w - - 4 29")
	/*
		Starting position:
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

		Random positions:
		"rnb1k2r/pp3pp1/4pq1p/2p5/1bBPP3/2N2N2/PP3PPP/R2Q1RK1 b kq - 1 9"
		"2br2k1/6p1/1p2p2p/5pb1/N2p4/PR6/1P3PPP/R5K1 b - - 4 29"

		Castling tests:
		r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1 // can castle
		rN2k1Nr/8/8/8/8/8/8/RN2K1NR w KQkq - 0 1 // cant castle because attacked
		r3k2r/8/8/4R3/4r3/8/8/R3K2R w KQkq - 0 1 // cant castle because blocked

		Enpassant tests:
		rnbqkbnr/pp1p1ppp/8/2pPp3/8/8/PPP1PPPP/RNBQKBNR w KQkq c6 0 3 // white enpassant
		rnbqkbnr/pppp1ppp/8/8/3PpPP1/8/PPP1P2P/RNBQKBNR b KQkq f3 0 3 // black enpassant

		Promotion tests:
		1pbr2k1/P5p1/4p2p/5pb1/N7/1R6/3p1PPP/R1P3K1 w - - 4 29 // white capture promotion
		2br2k1/6p1/1p2p2p/5pb1/N7/PR6/3p1PPP/R1P3K1 b - - 4 29 // black capture promotion
	*/
}

func (cb *ChessBoard) Test() {
	cb.PrintChessBoard()
	var moveList []Move
	cb.GenerateMoves(&moveList)
	PrintMoveList(moveList[:])
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
				fmt.Print("|.")
			}
		}
		fmt.Print("|\n")
	}
	fmt.Print("  a b c d e f g h\n")
}

func (cb ChessBoard) PrintChessBoard() {
	PrintBitboard(cb.WhitePieces | cb.BlackPieces)
	if SideToMove == White {
		fmt.Printf("\nSide to move: %s", "White")
	} else {
		fmt.Printf("\nSide to move: %s", "Black")
	}
	fmt.Printf("\nCastling rights (KQkq): ")
	if CastleRights&WhiteKingSide != 0 {
		fmt.Printf("K")
	}
	if CastleRights&WhiteQueenSide != 0 {
		fmt.Printf("Q")
	}
	if CastleRights&BlackKingSide != 0 {
		fmt.Printf("k")
	}
	if CastleRights&BlackQueenSide != 0 {
		fmt.Printf("q")
	}
	if Enpassant == -1 {
		fmt.Printf("\nEnpassant: %s", "-")
	} else {
		fmt.Printf("\nEnpassant: %s", IntToSquare[Enpassant])
	}
	fmt.Printf("\nHalf move clock: %d", HalfMoveClock)
	fmt.Printf("\nFull move counter: %d\n\n", FullMoveCounter)
}

// Used to test whether the function IsAttackedBySide() returns the correct output
/*
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
*/
// Used to print the hexadecimal bitboard representation,
// which is used to get bit masks
/*
func PrintBitboardHex(bitboard uint64) {
	fmt.Printf("%s\n", fmt.Sprintf("0x%X", bitboard))
}
*/

func PrintMove(move Move) {
	fmt.Printf("%4s%s%-6s%-10s%-12s%04b%9d\n", "",
		IntToSquare[move.GetMoveStart()],
		IntToSquare[move.GetMoveEnd()],
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
