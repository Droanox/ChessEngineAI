package eval

import "github.com/Droanox/ChessEngineAI/src/board"

///////////////////////////////////////////////////////////////////
// Masks
///////////////////////////////////////////////////////////////////

// Bitboard masks for files on the chessboard, index by square
var FileMasks = [64]uint64{}

// Bitboard masks for the king
var KingFileMasks = [64]uint64{}

// Bitboard masks for ranks on the chessboard, index by square
var RankMasks = [64]uint64{}

// Bitboard masks for partial file masks, index by square
var PartialMasks = [2][64]uint64{}

// Bitboard masks for king zones, index by square
var KingZoneMasks = [2][64]uint64{}

// Isolated pawn masks, index by square
var EastIsolatedMasks = [64]uint64{}
var WestIsolatedMasks = [64]uint64{}

// Passed pawn masks, index by sidetomove and by square
var PassedMasks = [2][64]uint64{}

var KingForwardSquares = [2][64]uint64{}

///////////////////////////////////////////////////////////////////
// Pawn
///////////////////////////////////////////////////////////////////

// if there is a doubled pawn, subtract this value from the evaluation
var doublePawnPenaltyHalf = [2]int{6, 16}

// If there is an isolated pawn, subtract this value from the evaluation
var isolatedPawnPenaltyHalf = [2]int{8, 22}

// If there is a past pawn, add this value to the evaluation
var pastPawnBonus = [7]int{0, 3, 11, 26, 55, 103, 158}

//var pastPawnBonus = [7]int{0, 2, 12, 18, 58, 140, 240}

///////////////////////////////////////////////////////////////////
// Knight
///////////////////////////////////////////////////////////////////

// Mobility bonus for knights
var knightMobility = 5

// The less pawns there are, the less knights are worth
var knightPawnPenalty = 8

///////////////////////////////////////////////////////////////////
// Bishop
///////////////////////////////////////////////////////////////////

// Mobility bonus for bishops
var bishopMobility = 6

// if there is a bishop pair, add this value to the evaluation
var bishopPairBonus = 20

///////////////////////////////////////////////////////////////////
// Rooks
///////////////////////////////////////////////////////////////////

// if there is a semi-open file, add this value to the evaluation
var rookSemiOpenFile = 7

// if there is an open file, add this value to the evaluation
var rookOpenFile = 13

// count the number of queens and rooks on the same file, and award this bonus
var stackedPieceBonus = 5

// if there is a rook on the seventh rank, add this value to the evaluation
var rookOnSeventh = 6

// The less pawns there are, the more rooks are worth
var rookPawnBonus = 8

// Mobility bonus for rooks
var rookMobility = [2]int{2, 4}

// See where the pawns are for the other side
var rook7thRank = [2]uint64{
	0: board.Rank7On, 1: board.Rank2On,
}

///////////////////////////////////////////////////////////////////
// King
///////////////////////////////////////////////////////////////////

// if there is a semi-open file, subtract this value to the evaluation
// var kingSemiOpenFile = 10

// if there is an open file, subtract this value to the evaluation
// var kingOpenFile = 17

// Multiply pawn squares by this weight
var pawnMultiplier = 15

// Replace king wiht a queen, and see how many squares it sees, multiply by this value
var kingTropismPenaltyMultiplier = 3

// Table to evaluate king safety, from stockfish
// chessprogramming.org/King_Safety#Attack_Units
var SafetyTable = []int{
	0, 0, 1, 2, 3, 5, 7, 9, 12, 15,
	18, 22, 26, 30, 35, 39, 44, 50, 56, 62,
	68, 75, 82, 85, 89, 97, 105, 113, 122, 131,
	140, 150, 169, 180, 191, 202, 213, 225, 237, 248,
	260, 272, 283, 295, 307, 319, 330, 342, 354, 366,
	377, 389, 401, 412, 424, 436, 448, 459, 471, 483,
	494, 500, 500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500, 500, 500,
}

var Threat [2]int

var KingAttackingPieces [2]int

///////////////////////////////////////////////////////////////////
// Queen
///////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////
// Piece and square values
///////////////////////////////////////////////////////////////////

// Changes every time eval is called
var PawnValue = 82

// Values from http://www.talkchess.com/forum3/viewtopic.php?f=2&t=68311&start=19
var (
	PieceValuesMG = [6]int{82, 337, 365, 477, 1025, 3000}
	PieceValuesEG = [6]int{94, 281, 297, 512, 936, 3000}
)

var (
	gamephaseInc = [12]int{0, 1, 1, 2, 4, 0, 0, 1, 1, 2, 4, 0}
	tableMG      = [12][64]int{}
	tableEG      = [12][64]int{}
)

var pawnMG = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	98, 134, 61, 95, 68, 126, 34, -11,
	-6, 7, 26, 31, 65, 56, 25, -20,
	-14, 13, 6, 21, 23, 12, 17, -23,
	-27, -2, -5, 12, 17, 6, 10, -25,
	-26, -4, -4, -10, 3, 3, 33, -12,
	-35, -1, -20, -23, -15, 24, 38, -22,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var pawnEG = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	178, 173, 158, 134, 147, 132, 165, 187,
	94, 100, 85, 67, 56, 53, 82, 84,
	32, 24, 13, 5, -2, 4, 17, 17,
	13, 9, -3, -7, -7, -8, 3, -1,
	4, 7, -6, 1, 0, -5, -1, -8,
	13, 8, 8, 10, 13, 0, 2, -7,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var knightMG = [64]int{
	-167, -89, -34, -49, 61, -97, -15, -107,
	-73, -41, 72, 36, 23, 62, 7, -17,
	-47, 60, 37, 65, 84, 129, 73, 44,
	-9, 17, 19, 53, 37, 69, 18, 22,
	-13, 4, 16, 13, 28, 19, 21, -8,
	-23, -9, 12, 10, 19, 17, 25, -16,
	-29, -53, -12, -3, -1, 18, -14, -19,
	-105, -21, -58, -33, -17, -28, -19, -23,
}

var knightEG = [64]int{
	-58, -38, -13, -28, -31, -27, -63, -99,
	-25, -8, -25, -2, -9, -25, -24, -52,
	-24, -20, 10, 9, -1, -9, -19, -41,
	-17, 3, 22, 22, 22, 11, 8, -18,
	-18, -6, 16, 25, 16, 17, 4, -18,
	-23, -3, -1, 15, 10, -3, -20, -22,
	-42, -20, -10, -5, -2, -20, -23, -44,
	-29, -51, -23, -15, -22, -18, -50, -64,
}

var bishopMG = [64]int{
	-29, 4, -82, -37, -25, -42, 7, -8,
	-26, 16, -18, -13, 30, 59, 18, -47,
	-16, 37, 43, 40, 35, 50, 37, -2,
	-4, 5, 19, 50, 37, 37, 7, -2,
	-6, 13, 13, 26, 34, 12, 10, 4,
	0, 15, 15, 15, 14, 27, 18, 10,
	4, 15, 16, 0, 7, 21, 33, 1,
	-33, -3, -14, -21, -13, -12, -39, -21,
}

var bishopEG = [64]int{
	-14, -21, -11, -8, -7, -9, -17, -24,
	-8, -4, 7, -12, -3, -13, -4, -14,
	2, -8, 0, -1, -2, 6, 0, 4,
	-3, 9, 12, 9, 14, 10, 3, 2,
	-6, 3, 13, 19, 7, 10, -3, -9,
	-12, -3, 8, 10, 13, 3, -7, -15,
	-14, -18, -7, -1, 4, -9, -15, -27,
	-23, -9, -23, -5, -9, -16, -5, -17,
}

var rookMG = [64]int{
	32, 42, 32, 51, 63, 9, 31, 43,
	27, 32, 58, 62, 80, 67, 26, 44,
	-5, 19, 26, 36, 17, 45, 61, 16,
	-24, -11, 7, 26, 24, 35, -8, -20,
	-36, -26, -12, -1, 9, -7, 6, -23,
	-45, -25, -16, -17, 3, 0, -5, -33,
	-44, -16, -20, -9, -1, 11, -6, -71,
	-19, -13, 1, 17, 16, 7, -37, -26,
}

var rookEG = [64]int{
	13, 10, 18, 15, 12, 12, 8, 5,
	11, 13, 13, 11, -3, 3, 8, 3,
	7, 7, 7, 5, 4, -3, -5, -3,
	4, 3, 13, 1, 2, 1, -1, 2,
	3, 5, 8, 4, -5, -6, -8, -11,
	-4, 0, -5, -1, -7, -12, -8, -16,
	-6, -6, 0, 2, -9, -9, -11, -3,
	-9, 2, 3, -1, -5, -13, 4, -20,
}

var queenMG = [64]int{
	-28, 0, 29, 12, 59, 44, 43, 45,
	-24, -39, -5, 1, -16, 57, 28, 54,
	-13, -17, 7, 8, 29, 56, 47, 57,
	-27, -27, -16, -16, -1, 17, -2, 1,
	-9, -26, -9, -10, -2, -4, 3, -3,
	-14, 2, -11, -2, -5, 2, 14, 5,
	-35, -8, 11, 2, 8, 15, -3, 1,
	-1, -18, -9, 10, -15, -25, -31, -50,
}

var queenEG = [64]int{
	-9, 22, 22, 27, 27, 19, 10, 20,
	-17, 20, 32, 41, 58, 25, 30, 0,
	-20, 6, 9, 49, 47, 35, 19, 9,
	3, 22, 24, 45, 57, 40, 57, 36,
	-18, 28, 19, 47, 31, 34, 39, 23,
	-16, -27, 15, 6, 9, 17, 10, 5,
	-22, -23, -30, -16, -16, -23, -36, -32,
	-33, -28, -22, -43, -5, -32, -20, -41,
}

var kingMG = [64]int{
	-65, 23, 16, -15, -56, -34, 2, 13,
	29, -1, -20, -7, -8, -4, -38, -29,
	-9, 24, 2, -16, -20, 6, 22, -22,
	-17, -20, -12, -27, -30, -25, -14, -36,
	-49, -1, -27, -39, -46, -44, -33, -51,
	-14, -14, -22, -46, -44, -30, -15, -27,
	1, 7, -8, -64, -43, -16, 9, 8,
	-15, 36, 12, -54, 8, -28, 24, 14,
}

var kingEG = [64]int{
	-74, -35, -18, -18, -11, 15, 4, -17,
	-12, 17, 14, 17, 17, 38, 23, 11,
	10, 17, 23, 15, 20, 45, 44, 13,
	-8, 22, 24, 27, 26, 33, 26, 3,
	-18, -4, 21, 24, 27, 23, 9, -11,
	-19, -3, 11, 21, 23, 16, 7, -9,
	-27, -11, 4, 13, 14, 4, -5, -17,
	-53, -34, -21, -11, -28, -14, -24, -43,
}

var piecesMG = [6][64]int{
	pawnMG,
	knightMG,
	bishopMG,
	rookMG,
	queenMG,
	kingMG,
}

var piecesEG = [6][64]int{
	pawnEG,
	knightEG,
	bishopEG,
	rookEG,
	queenEG,
	kingEG,
}

///////////////////////////////////////////////////////////////////
// General util
///////////////////////////////////////////////////////////////////

var pieceToColour = []int{
	0: board.White, 1: board.White, 2: board.White, 3: board.White, 4: board.White, 5: board.White,
	6: board.Black, 7: board.Black, 8: board.Black, 9: board.Black, 10: board.Black, 11: board.Black,
}

var offsetBySide = []int{
	board.White: +8, board.Black: -8,
}
