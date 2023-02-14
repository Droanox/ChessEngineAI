package engine

import (
	"github.com/Droanox/ChessEngineAI/src/board"
)

// centipawn values corresponding to pieces
// pawn, knight, bishop, rook, queen, king in that order
// with the first row as white and second as black

var (
	gamephaseInc = [12]int{0, 0, 1, 1, 1, 1, 2, 2, 4, 4, 0, 0}
	tableMG      = [12][64]int{}
	tableEG      = [12][64]int{}
)

func initValues() {
	for piece := board.Pawn - 1; piece < board.King; piece++ {
		for square := 0; square < 64; square++ {
			tableMG[piece][square] = pieceValuesMG[piece] + piecesMG[piece][square^56]
			tableEG[piece][square] = pieceValuesEG[piece] + piecesEG[piece][square^56]
			tableMG[piece+6][square] = pieceValuesMG[piece] + piecesMG[piece][square]
			tableEG[piece+6][square] = pieceValuesEG[piece] + piecesEG[piece][square]
		}
	}
}

func Init(board *board.ChessBoard) {
	// cb = *board

	initValues()
}

func Eval(cb board.ChessBoard) int {
	var pieceArr = []*uint64{
		0: &cb.WhitePawns, 1: &cb.WhiteKnights, 2: &cb.WhiteBishops, 3: &cb.WhiteRooks, 4: &cb.WhiteQueen, 5: &cb.WhiteKing,
		6: &cb.BlackPawns, 7: &cb.BlackKnights, 8: &cb.BlackBishops, 9: &cb.BlackRooks, 10: &cb.BlackQueen, 11: &cb.BlackKing,
	}

	var side, square, gamephase int
	var mg, eg [2]int

	for i, pieceBoard := range pieceArr {
		side = pieceToColour[i]
		for bitboard := *pieceBoard; bitboard != board.EmptyBoard; bitboard &= bitboard - 1 {
			square = board.BitScanForward(bitboard)
			mg[side] += tableMG[i][square]
			eg[side] += tableEG[i][square]
			gamephase += gamephaseInc[i]
		}
	}

	var scoreMG int = mg[board.SideToMove] - mg[1-board.SideToMove]
	var scoreEG int = eg[board.SideToMove] - eg[1-board.SideToMove]
	var phaseMG int = gamephase
	if phaseMG > 24 {
		phaseMG = 24
	}
	var phaseEG int = 24 - phaseMG
	return ((scoreMG * phaseMG) + (scoreEG * phaseEG)) / 24
}
