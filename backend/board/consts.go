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
	White = iota
	Black
)

const EmptyBoard uint64 = 0

const InitialPositionFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
