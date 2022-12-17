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
	WhiteQueensNum
	WhiteKingNum

	BlackPawnsNum
	BlackKnightsNum
	BlackBishopsNum
	BlackRooksNum
	BlackQueensNum
	BlackKingNum
)

const (
	White uint8 = iota
	Black
)

const EmptyBoard uint64 = 0

const initialPositionFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
