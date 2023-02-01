package board

import (
	"fmt"
	"time"
)

func (cb ChessBoard) PerftTest() {
	depth := 5

	start := time.Now()

	var moveList = []Move{}
	cb.GenerateMoves(&moveList)

	for i := 0; i < len(moveList); i++ {
		if !cb.MakeMove(moveList[i]) {
			continue
		}
		var cumNodes int64 = nodes
		cb.perftDriver(depth - 1)
		var oldNodes int64 = nodes - cumNodes
		cb.MakeBoard()
		fmt.Printf("%s%s:", IntToSquare[moveList[i].GetMoveStart()], IntToSquare[moveList[i].GetMoveEnd()])
		fmt.Printf(" %d\n", oldNodes)
	}

	elapsed := time.Since(start)
	fmt.Println("\nNodes searched:", nodes)
	fmt.Println("Time elapsed:", elapsed)
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

	for i := start; i < end; i++ {
		MagicInitWithReturn(int64(i))
		cb.CopyBoard()
		cb.Init("3k4/1r4r1/2b2b2/8/8/2B2B2/8/R3K2R w - - 0 1")

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
		cb.MakeBoard()
	}
	fmt.Println("\nBest seed found: ", bestSeed, " With best time: ", bestTime)
	MagicInit(int64(bestSeed))
}
