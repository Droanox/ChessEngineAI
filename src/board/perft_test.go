package board

import (
	"fmt"
	"testing"
	"time"
)

// Perft tests taken from Chess Programming wiki
// https://www.chessprogramming.org/Perft_Results
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
		nodes: 89941195,
	}, {
		Name:  "Wiki Position 6",
		FEN:   "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
		depth: 5,
		nodes: 164075551,
	},
}

// TestPerft is similar to PerftTest but runs through a list of tests and has a timer.
func TestPerft(t *testing.T) {
	cb := ChessBoard{}
	cb.Init()
	for _, test := range perftTests {
		start := time.Now()
		cb.CopyBoard()

		cb.ParseFen(test.FEN)
		cb.perftDriver(test.depth)
		elapsed := time.Since(start)
		fmt.Printf("%s\nFen: %s\nNodes searched: %d\nTime elapsed: %v\n", test.Name, test.FEN, nodes, elapsed)
		if nodes != test.nodes {
			fmt.Printf("\u2718\n\n")
			t.Errorf("Incorrect node count\nReceived: %d\nNeeded: %d", nodes, test.nodes)
		} else {
			fmt.Printf("\u2714\n\n")
		}
		nodes = 0

		cb.MakeBoard()
	}
	fmt.Println()
}
