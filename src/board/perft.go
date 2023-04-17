package board

import (
	"fmt"
	"time"
)

// PerftTest runs a series of perft tests to verify that the move generation is working correctly
func PerftTest(fen string, depth int) {
	var totalNodes int64 = 0

	cb := ChessBoard{}
	cb.Init()
	cb.ParseFen(fen)

	var moveList = []Move{}
	cb.GenerateMoves(&moveList)

	for i := 0; i < len(moveList); i++ {
		if !cb.MakeMove(moveList[i]) {
			continue
		}
		cb.perftDriver(depth - 1)
		fmt.Print(IndexToSquare[moveList[i].GetMoveStart()]+IndexToSquare[moveList[i].GetMoveEnd()], ": ", nodes, "\n")
		totalNodes += nodes
		nodes = 0
		cb.MakeBoard()
	}
	fmt.Println("\nTotal nodes:", totalNodes)
	/*
		for _, test := range PerftTests {
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
	*/
}

// PerftTestEXTRA is similar to PerftTest but runs through a list of tests and has a timer.
func PerftTestEXTRA() {
	cb := ChessBoard{}
	cb.Init()
	for i, test := range PerftTests {
		start := time.Now()
		cb.CopyBoard()

		cb.ParseFen(test.FEN)
		cb.perftDriver(test.depth)
		elapsed := time.Since(start)
		fmt.Print(test.Name, "\nFen: ", test.FEN, "\nNodes searched: ", nodes, "\nTime elapsed: ", elapsed, "\n\n")
		if nodes != test.nodes {
			fmt.Print("ERROR!\nReceived: ", nodes, "\nNeeded:   ", test.nodes)
			// New line after every test, stop before the last test
			if i < len(PerftTests)-1 {
				fmt.Print("\n\n")
			}
		}
		nodes = 0

		cb.MakeBoard()
	}
	fmt.Println()
}

// PerftDriver is a recursive function that counts the number of nodes at a given depth
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

// PerftTestFindMagic is a function that finds the best seed for the magic bitboards
// It is not used in the program, but it is useful for finding the best seed
// Not fully accurate, but it is a good starting point
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
