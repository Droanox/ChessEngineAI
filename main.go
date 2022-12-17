package main

import (
	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/graphics"
)

/*import (
	"encoding/binary"
	"fmt"
)

const (
	FileAMask uint64 = 0xFF << (8 * iota)
	FileBMask
	FileCMask
	FileDMask
	FileEMask
	FileFMask
	FileGMask
	FileHMask
) */

func main() {
	cb := board.ChessBoard{}
	cb.Init()
	graphics.Run(cb)

	//bitboard = board.MaskPawnAttacks(4, board.Black)
	//board.PrintBitboard(bitboard)
	//fmt.Print(144680345676153346 << 1)
}
