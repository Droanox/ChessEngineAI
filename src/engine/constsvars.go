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

// killerMoves is used to store the killer moves
// The first index is the number of killers (2)
// The second index is the ply
// The value is the move
var killerMoves = [2][board.MaxPly]board.Move{}

// historyMoves is used to store the history of moves
// The first index is the number of pieces (12)
// The second index is the ply
// The value is the move
var historyMoves = [12][board.MaxPly]int{}

///////////////////////////////////////////////////////////////////
// Principal variation
///////////////////////////////////////////////////////////////////

// Principal variation length
// The index is the ply
// The value is the length of the principal variation
var pvLength [board.MaxPly]int

// Principal variation table
// The first index is the ply
// The second index is the length of the principal variation
// The value is the move
var pvTable [board.MaxPly][board.MaxPly]board.Move

// pvFollowed is used to determine if the principal variation was followed
// if true then the principal variation was followed
// if false then the principal variation was not followed
var pvFollowed bool

///////////////////////////////////////////////////////////////////
// Late move reduction
///////////////////////////////////////////////////////////////////

// fullDepthMoves is the number of moves searched before
// late move reduction is used
var fullDepthMoves int = 5 // Changeable by user

// reductionLimit is the maximum number of reductions
// that can be performed
var reductionLimit int = 3 // Changeable by user

///////////////////////////////////////////////////////////////////
// Null move pruning
///////////////////////////////////////////////////////////////////

// R is the depth reduction factor for null move pruning
// This is the number of plies to reduce the depth by
var nullMoveReduction int = 2 // Changeable by user

// nullMoveDepth is the depth at which null move pruning is used
var nullMoveDepth int = 3 // Changeable by user

///////////////////////////////////////////////////////////////////
// General util
///////////////////////////////////////////////////////////////////

// moveOrderOffset is used to offset the move ordering score
const moveOrderOffset = 1000 //math.MaxInt32 - 1024

// nodes is the number of nodes searched
// This is used to print the number of nodes searched
// in the search info
// This is reset every time a search is performed
// and is incremented every time a node is searched
var nodes int
