package search

import (
	"github.com/Droanox/ChessEngineAI/src/board"
)

///////////////////////////////////////////////////////////////////
// Move ordering
///////////////////////////////////////////////////////////////////

// Indexed by MVV_LVA[Victim][Attacker]
// Use the board package as the index will always be no. pieces + 1,
// since there is an empty piece initialized into this array
var MVV_LVA = [7][7]int{
	{0, 00, 00, 00, 00, 00, 00},
	{0, 15, 14, 13, 12, 11, 10},
	{0, 25, 24, 23, 22, 21, 20},
	{0, 35, 34, 33, 32, 31, 30},
	{0, 45, 44, 43, 42, 41, 40},
	{0, 55, 54, 53, 52, 51, 50},
	{0, 65, 64, 63, 62, 61, 60},
}

var killerMoves = [2][board.MaxPly]board.Move{}

var historyMoves = [12][board.MaxPly]int{}

///////////////////////////////////////////////////////////////////
// General util
///////////////////////////////////////////////////////////////////

const moveOrderOffset = 1000 //math.MaxInt32 - 1024

var nodes int

var bestMove board.Move
