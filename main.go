package main

import (
	"fmt"

	"github.com/Droanox/ChessEngineAI/backend/board"
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
	/*s := make([]byte, 8)
	b := board.NewChessBoard()
	binary.LittleEndian.PutUint64(s, 0x03f79d71b4cb0a89)
	fmt.Printf("%08b", s)
	board.PrintOutSimple(b) */
	var bitboard uint64
	board.SetBit(&bitboard, 4)
	board.PrintBitboard(bitboard)
	bitboard = board.MaskPawnAttacks(4, board.White)
	board.PrintBitboard(bitboard)
	fmt.Print("\n")
}
