package main

import (
	"fmt"

	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/engine"
	"github.com/Droanox/ChessEngineAI/src/eval"
	"github.com/Droanox/ChessEngineAI/src/uci"
)

// command to run to create the .exe
// fyne package -os windows -icon ChessEngineAI.png
func main() {
	debug := true

	if debug {
		cb := board.ChessBoard{}
		cb.Init()
		eval.Init()
		engine.Init()
		cb.ParseFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/r3K2R w K-kq -")
		engine.HashKey = engine.GenHash(cb)
		cb.PrintChessBoard()
		fmt.Printf("%0x", engine.HashKey)
		//engine.Search(12, &cb)
	} else {
		uci.Run()
	}
}

/*
	PERFT tests:
	"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -"

	PERFT test find magic:
	"3k4/1r4r1/2b2b2/8/8/2B2B2/8/R3K2R w - - 0 1"

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
	rnbqkbnr/pp1ppppp/8/2pP4/8/8/PPP1PPPP/RNBQKBNR b KQkq - 0 2 // black enpassant is set

	Promotion tests:
	1pbr2k1/P5p1/4p2p/5pb1/N7/1R6/3p1PPP/R1P3K1 w - - 4 29 // white capture promotion
	2br2k1/6p1/1p2p2p/5pb1/N7/PR6/3p1PPP/R1P3K1 b - - 4 29 // black capture promotion
*/
