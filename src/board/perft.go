package board

import (
	"fmt"
	"time"
)

func (cb ChessBoard) PerftTest() {
	cb.Init()
	for _, test := range perftTests {
		cb.CopyBoard()

		cb.ParseFen(test.FEN)
		cb.perftDriver(test.depth)
		fmt.Println(test.Name, "\nNodes searched:", nodes)
		if nodes != test.nodes {
			fmt.Println("ERROR!")
		}
		nodes = 0

		cb.MakeBoard()
	}
}

func (cb ChessBoard) PerftTestTimer() {
	cb.Init()
	for _, test := range perftTests {
		start := time.Now()
		cb.CopyBoard()

		cb.ParseFen(test.FEN)
		cb.perftDriver(test.depth)
		elapsed := time.Since(start)
		fmt.Println(test.Name, "\nNodes searched:", nodes, "\nTime elapsed:", elapsed)
		if nodes != test.nodes {
			fmt.Println("ERROR!", nodes, "!=", test.nodes)
		}
		nodes = 0

		cb.MakeBoard()
	}
	fmt.Println()
}

func (cb ChessBoard) perftDriver(depth int) {
	if depth == 0 {
		nodes++
		return
	}

	var moveList = []Move{}
	cb.GenerateMoves(&moveList)

	for i := 0; i < len(moveList); i++ {
		if !cb.MakeMove(moveList[i]) {
			continue
		}
		cb.perftDriver(depth - 1)
		cb.MakeBoard()
	}
}

func (cb ChessBoard) PerftTestFindMagic(start int, end int) {
	bestSeed := -1
	var bestTime time.Duration = 100000000000
	cb.Init()

	for i := start; i < end; i++ {
		MagicInitWithReturn(int64(i))
		cb.ParseFen("3k4/1r4r1/2b2b2/8/8/2B2B2/8/R3K2R w - - 0 1")

		depth := 5

		start := time.Now()

		var moveList = []Move{}
		cb.GenerateMoves(&moveList)

		for i := 0; i < len(moveList); i++ {
			if !cb.MakeMove(moveList[i]) {
				continue
			}
			cb.perftDriver(depth - 1)
			cb.MakeBoard()
		}

		elapsed := time.Since(start)
		fmt.Println("\nNodes searched:", nodes)
		fmt.Println("Seed used:", i)
		fmt.Println("Time elapsed:", elapsed)
		nodes = 0

		if time.Duration(elapsed.Milliseconds()) < time.Duration(bestTime.Milliseconds()) {
			bestSeed = i
			bestTime = elapsed
			fmt.Printf("New best Seed: %d New best Time: ", bestSeed)
			fmt.Println(bestTime)
		}
	}
	fmt.Println("\nBest seed found: ", bestSeed, " With best time: ", bestTime)
	MagicInit(int64(bestSeed))
}
