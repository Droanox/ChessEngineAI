package board

import (
	"fmt"
	"time"
)

func (cb ChessBoard) PerftTest() {
	depth := 6

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
