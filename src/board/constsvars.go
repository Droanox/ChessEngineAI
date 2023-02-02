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
	// Ply is incremented after each coard copy, and decremented after each board make (paste)
	Ply int = -1
)

var IntToPiece = [7]string{
	"Empty", "Pawn", "Knight", "Bishop", "Rook", "Queen", "King",
}

var PromotionToPiece = map[int]int{
	MoveKnightPromotion: 2, MoveBishopPromotion: 3, MoveRookPromotion: 4, MoveQueenPromotion: 5,
	MoveKnightPromotionCapture: 2, MoveBishopPromotionCapture: 3, MoveRookPromotionCapture: 4, MoveQueenPromotionCapture: 5,
}

var SideToOffset = []int{
	White: -8, Black: +8,
}

// Used to check if a moved piece affects the current castlerights var, if so,
// update it accordingly
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

var perftTests = []struct {
	Name  string
	FEN   string
	depth int
	nodes int64
}{
	{
		Name:  "Initial Postion",
		FEN:   "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		depth: 6,
		nodes: 119060324,
	}, {
		Name:  "Wiki Position 2",
		FEN:   "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -",
		depth: 5,
		nodes: 193690690,
	}, {
		Name:  "Wiki Position 3",
		FEN:   "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - -",
		depth: 6,
		nodes: 11030083,
	}, {
		Name:  "Wiki Position 4.1",
		FEN:   "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		depth: 5,
		nodes: 15833292,
	}, {
		Name:  "Wiki Position 4.2",
		FEN:   "r2q1rk1/pP1p2pp/Q4n2/bbp1p3/Np6/1B3NBn/pPPP1PPP/R3K2R b KQ - 0 1",
		depth: 5,
		nodes: 15833292,
	}, {
		Name:  "Wiki Position 5",
		FEN:   "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		depth: 5,
		nodes: 89941194,
	}, {
		Name:  "Wiki Position 6",
		FEN:   "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
		depth: 5,
		nodes: 164075551,
	},
}

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

var bishopBits = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

var rookBits = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

// Used to change the square index to the x-y mapping of it
var IndexToSquare = [65]string{
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
var SquareToIndex = map[string]int{
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

const (
	Empty int = iota
	WhitePawns
	WhiteKnights
	WhiteBishops
	WhiteRooks
	WhiteKing
	WhiteQueen

	BlackPawns
	BlackKnights
	BlackBishops
	BlackRooks
	BlackKing
	BlackQueen
)

var PieceToASCII = []string{
	0: ".",
	1: "♙", 2: "♘", 3: "♗", 4: "♖", 5: "♔", 6: "♕",
	7: "♟", 8: "♞", 9: "♝", 10: "♜", 11: "♚", 12: "♛",
}

///////////////////////////////////////////////////////////////////
// Magic number util consts and vars
///////////////////////////////////////////////////////////////////

// Best magic number seed found so far: 15

var rookMagicNumber = [64]uint64{
	1188950372496965666, 1170936040556331072, 648531610290888704, 5044036049991417992, 1585284669612508168, 144119622663143696, 288230930773182480, 4755801481985165586,
	1189372515164895872, 38421471764152448, 2308235615243608064, 563027264932928, 2523423187686918144, 2306405976615550984, 10674376643465838596, 2378463558642205193,
	90213279793692672, 75041670696961, 576602039816036352, 4504149517406336, 1442278884473047040, 37718748672819456, 144679241836134912, 10995393642756,
	2684145517499269248, 36312496792961056, 36436720132071425, 2289221864784128, 6756503248179204, 761671304226021392, 9259968186185417216, 2305860919229679753,
	141012374650914, 141012374659072, 301741346834613249, 20301438546612225, 5226427384768497792, 288232577330840576, 3531387256851988484, 576742915582201604,
	117234466312716288, 5769112501883125760, 144132780798804096, 2452491609787400200, 2533309151248388, 8444387008971776, 576462951393853696, 145258972505571329,
	4616260177925046528, 2306485682366851584, 619316443076362496, 20266234898546816, 659289165065715840, 306385520739745920, 9223380841555559424, 72339077612978432,
	9259436301730480146, 288270242340208770, 158468220749954, 576470648445210625, 9252927201503881217, 577023708837982210, 2449994484127629860, 146649567232000386,
}

var bishopMagicNumber = [64]uint64{
	290517652400111648, 9802093389561212928, 1162511462273384448, 2287534547664898, 1157574778073124912, 9296012381182246928, 324823291398225922, 142940817081408,
	869212328892566024, 176207737830834704, 18159130125313, 19144713785626624, 2387193693246333216, 9016143256748354, 4684317566393518786, 2522015937390075904,
	2251869106227202, 4909205103304573000, 9228441129125758210, 1139265882818560, 289356542381195266, 4612847107010215936, 5620773861497250352, 288266668659048960,
	325402674392531204, 585768290371242000, 3461029507908256020, 1170940301167100065, 18295944369946624, 3386771774119936, 563224865407536, 142386791715840,
	5084160155787264, 5765258451114657824, 585503342105986176, 1152925919833752080, 4611968610095726624, 11538260733599025176, 2594381545108799748, 4653556549732663809,
	36314679710064771, 144406077690621952, 9512167570724620801, 10377560454305189920, 79371265114368, 2458969865491316768, 2333014142791256192, 292808813438894212,
	5206732923911278592, 1153062827302332416, 18859924042809352, 2324561104862741120, 9228016618525950340, 22526931803079680, 5647153061101568, 9009003149427216,
	11547248694889963584, 76968030921232, 722854132909641728, 4719772478205888514, 578727946423915520, 19140315865678404, 18021065104098562, 2329636476277227584,
}
