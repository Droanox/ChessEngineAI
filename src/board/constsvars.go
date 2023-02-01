package board

///////////////////////////////////////////////////////////////////
// Masks
///////////////////////////////////////////////////////////////////
// Hexadecimal values for constants taken from
// chess programming wiki square mapping considerations
const (
	FileAOn uint64 = 0x0101010101010101 << iota
	FileBOn
	FileCOn
	FileDOn
	FileEOn
	FileFOn
	FileGOn
	FileHOn
)

const (
	Rank1On uint64 = 0x00000000000000FF << (8 * iota)
	Rank2On
	Rank3On
	Rank4On
	Rank5On
	Rank6On
	Rank7On
	Rank8On
)

const (
	FileABOn uint64 = 0x303030303030303
	FileGHOn uint64 = 0xc0c0c0c0c0c0c0c0
)

const (
	White = 0b0
	Black = 0b1
)

const EmptyBoard uint64 = 0x0000000000000000
const UniverseBoard uint64 = 0xffffffffffffffff

///////////////////////////////////////////////////////////////////
// FEN string consts and vars
///////////////////////////////////////////////////////////////////
const (
	WhiteKingSide  = 0b0001
	WhiteQueenSide = 0b0010
	BlackKingSide  = 0b0100
	BlackQueenSide = 0b1000
)

var AllPieceNames = []string{
	"WhitePawns", "WhiteKnights", "WhiteBishops", "WhiteRooks", "WhiteQueen", "WhiteKing",
	"BlackPawns", "BlackKnights", "BlackBishops", "BlackRooks", "BlackQueen", "BlackKing",
}

var SideToMoveMap = map[string]int{
	"w": White, "b": Black,
}
var CastleMap = map[rune]int{
	'K': WhiteKingSide, 'Q': WhiteQueenSide, 'k': BlackKingSide, 'q': BlackQueenSide,
}

// The standard starting position for a Chess Board in FEN notation
const InitialPositionFen string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var (
	// White = 0, Black = 1
	SideToMove int = 0

	// None = 0, white king side = 1, white queen side = 2, black king side = 4, black queen side = 8,
	// and all available will be the additions of them e.g.
	// 13 = 1101 = white king side, black king side, black queen side
	CastleRights int = 0

	// -1 for no enpassant, and 0-63 for square index of enpassant
	Enpassant int = -1

	// Counts every half move made, resets to 0 if a pawn is moved or piece is taken,
	// used for the 50 move draw rule
	HalfMoveClock int = 0

	// Counts the number of full moves made in a game, it starts at 1 and increments
	// every time black moves
	FullMoveCounter int = 0
)

///////////////////////////////////////////////////////////////////
// Movegen consts and vars
///////////////////////////////////////////////////////////////////

const (
	EmptyPiece = 0x0 // 0000 0000 0000 0000 0000 0000
	Pawn       = 0x1 // 0000 0000 0001 0000 0000 0000
	Knight     = 0x2 // 0000 0000 0010 0000 0000 0000
	Bishop     = 0x3 // 0000 0000 0011 0000 0000 0000
	Rook       = 0x4 // 0000 0000 0100 0000 0000 0000
	Queen      = 0x5 // 0000 0000 0101 0000 0000 0000
	King       = 0x6 // 0000 0000 0110 0000 0000 0000
)

const (
	/*
		MoveStart                  = 0x3f   //  0000 0000 0000 0000 0011 1111
		MoveEnd                    = 0xfc0  //  0000 0000 0000 1111 1100 0000
	*/
	MoveQuiet                  = 0x0 // 0000 0000 0000 0000 0000 0000
	MoveDoublePawn             = 0x1 // 0001 0000 0000 0000 0000 0000
	MoveKingCastle             = 0x2 // 0010 0000 0000 0000 0000 0000
	MoveQueenCastle            = 0x3 // 0011 0000 0000 0000 0000 0000
	MoveCaptures               = 0x4 // 0100 0000 0000 0000 0000 0000
	MoveEnpassantCapture       = 0x5 // 0101 0000 0000 0000 0000 0000
	MoveKnightPromotion        = 0x8 // 1000 0000 0000 0000 0000 0000
	MoveBishopPromotion        = 0x9 // 1001 0000 0000 0000 0000 0000
	MoveRookPromotion          = 0xA // 1010 0000 0000 0000 0000 0000
	MoveQueenPromotion         = 0xB // 1011 0000 0000 0000 0000 0000
	MoveKnightPromotionCapture = 0xC // 1100 0000 0000 0000 0000 0000
	MoveBishopPromotionCapture = 0xD // 1101 0000 0000 0000 0000 0000
	MoveRookPromotionCapture   = 0xE // 1110 0000 0000 0000 0000 0000
	MoveQueenPromotionCapture  = 0xF // 1111 0000 0000 0000 0000 0000
)

var (
	ChessBoardCopies [100]ChessBoard
	AspectsCopies    [100][5]int
	Ply              int = -1
)

var IntToPiece = [7]string{
	"Empty", "Pawn", "Knight", "Bishop", "Rook", "Queen", "King",
}

var PromotionToPiece = map[int]int{
	MoveKnightPromotion: 2, MoveBishopPromotion: 3, MoveRookPromotion: 4, MoveQueenPromotion: 5,
	MoveKnightPromotionCapture: 2, MoveBishopPromotionCapture: 3, MoveRookPromotionCapture: 4, MoveQueenPromotionCapture: 5,
}

var SideToOffset = map[int]int{
	White: -8, Black: +8,
}

// Used to check if a moved piece affects the current castlerights var, if so,
// update it accordingly
/*
var CastleRightsUpdate = [64]int{
	7, 15, 15, 15, 3, 15, 15, 11,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	13, 15, 15, 15, 12, 15, 15, 14,
}
*/
var CastleRightsUpdate = [64]int{
	13, 15, 15, 15, 12, 15, 15, 14,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	7, 15, 15, 15, 3, 15, 15, 11,
}

///////////////////////////////////////////////////////////////////
// Perft consts and vars
///////////////////////////////////////////////////////////////////

var nodes int64

///////////////////////////////////////////////////////////////////
// General util consts and vars
///////////////////////////////////////////////////////////////////
// Kim Walisch's proposed ones' decrement to compute
// the least significant 1 bit used in BitScanForward()
var index64 = [64]int{
	0, 47, 1, 56, 48, 27, 2, 60,
	57, 49, 41, 37, 28, 16, 3, 61,
	54, 58, 35, 52, 50, 42, 21, 44,
	38, 32, 29, 23, 17, 11, 4, 62,
	46, 55, 26, 59, 40, 36, 15, 53,
	34, 51, 20, 43, 31, 22, 10, 45,
	25, 39, 14, 33, 19, 30, 9, 24,
	13, 18, 8, 12, 7, 6, 5, 63,
}

// Used to change the square index to the x-y mapping of it
var IntToSquare = [65]string{
	"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
	"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
	"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
	"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
	"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
	"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
	"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
	"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
	"-",
}

// Used to change a x-y mapping of a chessboard to the squares index.
var SquareToInt = map[string]int{
	"a1": 0, "b1": 1, "c1": 2, "d1": 3, "e1": 4, "f1": 5, "g1": 6, "h1": 7,
	"a2": 8, "b2": 9, "c2": 10, "d2": 11, "e2": 12, "f2": 13, "g2": 14, "h2": 15,
	"a3": 16, "b3": 17, "c3": 18, "d3": 19, "e3": 20, "f3": 21, "g3": 22, "h3": 23,
	"a4": 24, "b4": 25, "c4": 26, "d4": 27, "e4": 28, "f4": 29, "g4": 30, "h4": 31,
	"a5": 32, "b5": 33, "c5": 34, "d5": 35, "e5": 36, "f5": 37, "g5": 38, "h5": 39,
	"a6": 40, "b6": 41, "c6": 42, "d6": 43, "e6": 44, "f6": 45, "g6": 46, "h6": 47,
	"a7": 48, "b7": 49, "c7": 50, "d7": 51, "e7": 52, "f7": 53, "g7": 54, "h7": 55,
	"a8": 56, "b8": 57, "c8": 58, "d8": 59, "e8": 60, "f8": 61, "g8": 62, "h8": 63,
	"-": 64,
}
