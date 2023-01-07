package board

const (
	FileAOn uint64 = 72340172838076673 << iota
	FileBOn
	FileCOn
	FileDOn
	FileEOn
	FileFOn
	FileGOn
	FileHOn
)

const (
	Rank1On uint64 = 255 << (8 * iota)
	Rank2On
	Rank3On
	Rank4On
	Rank5On
	Rank6On
	Rank7On
	Rank8On
)

const (
	WhitePawnsNum = iota
	WhiteKnightsNum
	WhiteBishopsNum
	WhiteRooksNum
	WhiteQueenNum
	WhiteKingNum

	BlackPawnsNum
	BlackKnightsNum
	BlackBishopsNum
	BlackRooksNum
	BlackQueenNum
	BlackKingNum
)

const (
	White uint8 = iota
	Black
)

const EmptyBoard uint64 = 0x0000000000000000
const UniverseBoard uint64 = 0xffffffffffffffff
const initialPositionFen string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var AllPieceNames = []string{
	"WhitePawns", "WhiteKnights", "WhiteBishops", "WhiteRooks", "WhiteQueen", "WhiteKing",
	"BlackPawns", "BlackKnights", "BlackBishops", "BlackRooks", "BlackQueen", "BlackKing"}

// Kim Walisch's proposed ones' decrement to compute
// the least significant 1 bit
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
