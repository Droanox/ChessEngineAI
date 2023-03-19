package engine

import (
	"time"

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
var killerMoves [2][board.MaxPly]board.Move

// historyMoves is used to store the history of moves
// The first index is the number of pieces (12)
// The second index is the ply
// The value is the move
var historyMoves [12][board.MaxPly]int

///////////////////////////////////////////////////////////////////
// Aspiration windows
///////////////////////////////////////////////////////////////////

// aspirationWindow is the size of the aspiration window used in iterative deepening
const aspirationWindow int = 50 // Changeable by user

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
const fullDepthMoves int = 4 // Changeable by user

// reductionLimit is the maximum number of reductions
// that can be performed
const reductionLimit int = 3 // Changeable by user

///////////////////////////////////////////////////////////////////
// Null move pruning
///////////////////////////////////////////////////////////////////

// R is the depth reduction factor for null move pruning
// This is the number of plies to reduce the depth by
const nullMoveReduction int = 2 // Changeable by user

// nullMoveDepth is the depth at which null move pruning is used
const nullMoveDepth int = 3 // Changeable by user

///////////////////////////////////////////////////////////////////
// Transposition table
///////////////////////////////////////////////////////////////////

// noHash is used to determine if the hash table entry is empty
const noHash int = maxScore + 1

// hashFlagExact is used to determine if the hash table entry is exact
const hashFlagExact int = 0

// hashFlagAlpha is used to determine if the hash table entry is alpha
const hashFlagAlpha int = 1

// hashFlagBeta is used to determine if the hash table entry is beta
const hashFlagBeta int = 2

// hashSize is the size of the hash table
const hashSize uint64 = 2 << 19 // default:16MB Changeable by user

// ttSize is the size of the transposition table
var tt [hashSize]TranspositionTable

///////////////////////////////////////////////////////////////////
// Time management
///////////////////////////////////////////////////////////////////

// TimeControl is used to determine if time control is used
var TimeControl bool = false

// start is the time the search started
var start time.Time

// stopTime is the time the search should stop
var StopTime int64

// isStopped is used to determine if the search should be stopped
var IsStopped bool

///////////////////////////////////////////////////////////////////
// General util
///////////////////////////////////////////////////////////////////

const maxScore = 100000

const minScore = -100000

const MateValue = minScore + 1000

const MateScore = MateValue + 1000

// moveOrderOffset is used to offset the move ordering score
const moveOrderOffset = 100

// nodes is the number of nodes searched
// This is used to print the number of nodes searched
// in the search info
// This is reset every time a search is performed
// and is incremented every time a node is searched
var nodes int
