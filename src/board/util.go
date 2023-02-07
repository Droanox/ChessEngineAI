package board

import (
	"strconv"
	"strings"
)

///////////////////////////////////////////////////////////////////
// General util
///////////////////////////////////////////////////////////////////

func setBit(bitboard *uint64, square int) {
	*bitboard |= (1 << uint64(square))
}

func getBit(bitboard uint64, square int) uint64 {
	return bitboard & (1 << uint64(square))
}

func popBit(bitboard *uint64, square int) {
	*bitboard ^= getBit(*bitboard, square)
}

func isBitOn(bitboard uint64, square int) bool {
	return bitboard == (bitboard | (1 << uint64(square)))
}

///////////////////////////////////////////////////////////////////
// Move generation
///////////////////////////////////////////////////////////////////

// The "Brian Kernighan's way" of counting bits on a bitboard,
// implementation idea from chess programming wiki
func BitCount(bitboard uint64) int {
	var count int
	for bitboard != EmptyBoard {
		count++
		bitboard &= bitboard - 1
	}
	return count
}

// Kim Walisch's proposed ones' decrement to compute
// the least significant 1 bit used in BitScanForward()
func BitScanForward(bitboard uint64) int {
	const debruijn64 uint64 = 0x03f79d71b4cb0a89
	if bitboard != 0 {
		return index64[((bitboard^(bitboard-1))*debruijn64)>>58]
	}
	return -1
}

func (cb ChessBoard) GetPieceType(square int) int {
	indexMask := indexMasks[square]
	switch {
	case (cb.WhitePawns|cb.BlackPawns)&indexMask != EmptyBoard:
		return Pawn
	case (cb.WhiteKnights|cb.BlackKnights)&indexMask != EmptyBoard:
		return Knight
	case (cb.WhiteBishops|cb.BlackBishops)&indexMask != EmptyBoard:
		return Bishop
	case (cb.WhiteRooks|cb.BlackRooks)&indexMask != EmptyBoard:
		return Rook
	case (cb.WhiteQueen|cb.BlackQueen)&indexMask != EmptyBoard:
		return Queen
	case (cb.WhiteKing|cb.BlackKing)&indexMask != EmptyBoard:
		return King
	}
	return EmptyPiece
}

///////////////////////////////////////////////////////////////////
// GUI
///////////////////////////////////////////////////////////////////

func (cb *ChessBoard) GetPieceString(square int) string {
	for _, p := range AllPieceNames {
		bitboard := cb.GetPiecesBitboardString(p)
		if isBitOn(bitboard, square) {
			return p
		}
	}
	return "Empty"
}

func (cb *ChessBoard) GetPiecesBitboardString(piece string) uint64 {
	pieceMap := map[string]uint64{
		"WhitePawns": cb.WhitePawns, "WhiteKnights": cb.WhiteKnights, "WhiteBishops": cb.WhiteBishops,
		"WhiteRooks": cb.WhiteRooks, "WhiteKing": cb.WhiteKing, "WhiteQueen": cb.WhiteQueen,

		"BlackPawns": cb.BlackPawns, "BlackKnights": cb.BlackKnights, "BlackBishops": cb.BlackBishops,
		"BlackRooks": cb.BlackRooks, "BlackKing": cb.BlackKing, "BlackQueen": cb.BlackQueen,
	}

	return pieceMap[piece]
}

func (cb *ChessBoard) GetPieceInt(square int) int {
	for i := 0; i <= BlackKing; i++ {
		bitboard := cb.GetPiecesBitboardInt(i)
		if isBitOn(bitboard, square) {
			return i
		}
	}
	return Empty
}

func (cb *ChessBoard) GetPiecesBitboardInt(piece int) uint64 {
	pieceMap := []uint64{
		Empty:      EmptyBoard,
		WhitePawns: cb.WhitePawns, WhiteKnights: cb.WhiteKnights, WhiteBishops: cb.WhiteBishops,
		WhiteRooks: cb.WhiteRooks, WhiteQueen: cb.WhiteQueen, WhiteKing: cb.WhiteKing,

		BlackPawns: cb.BlackPawns, BlackKnights: cb.BlackKnights, BlackBishops: cb.BlackBishops,
		BlackRooks: cb.BlackRooks, BlackQueen: cb.BlackQueen, BlackKing: cb.BlackKing,
	}

	return pieceMap[piece]
}

// Was used but made redundant, code is kept if it's ever needed again
/*
func (cb ChessBoard) Type(num int) string {
	b := reflect.TypeOf(cb)
	return b.Field(num).Name
}
*/

///////////////////////////////////////////////////////////////////
// Parse FEN
///////////////////////////////////////////////////////////////////

// Parses a fen string for example: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNN w KQkq - 0 1"
// onto the chessboard and maps every pieces bitboard to the relevant pieces

func (cb *ChessBoard) ParseFen(fen string) {
	ChessBoard{}.CopyBoard()
	cb.MakeBoard()
	fenRep := strings.Fields(fen)
	var file int
	var rank = 7

	pieceMap := map[rune]*uint64{
		'P': &cb.WhitePawns, 'N': &cb.WhiteKnights, 'B': &cb.WhiteBishops, 'R': &cb.WhiteRooks, 'K': &cb.WhiteKing, 'Q': &cb.WhiteQueen,
		'p': &cb.BlackPawns, 'n': &cb.BlackKnights, 'b': &cb.BlackBishops, 'r': &cb.BlackRooks, 'k': &cb.BlackKing, 'q': &cb.BlackQueen,
	}

	for _, val := range fenRep[0] {
		switch val {
		case '/':
			file = 0
			rank--
		case '1', '2', '3', '4', '5', '6', '7', '8':
			file += int(val - '0')
		default:
			setBit(pieceMap[val], (8*rank)+file)
			file++
		}
	}

	cb.WhitePieces = cb.WhitePawns | cb.WhiteKnights | cb.WhiteBishops | cb.WhiteRooks | cb.WhiteQueen | cb.WhiteKing
	cb.BlackPieces = cb.BlackPawns | cb.BlackKnights | cb.BlackBishops | cb.BlackRooks | cb.BlackQueen | cb.BlackKing

	SideToMove = SideToMoveMap[fenRep[1]]

	for _, val := range fenRep[2] {
		CastleRights += CastleMap[val]
	}

	Enpassant = SquareToIndex[fenRep[3]]

	if len(fenRep) > 4 {
		HalfMoveClock, _ = strconv.Atoi(fenRep[4])
		FullMoveCounter, _ = strconv.Atoi(fenRep[5])
	} else {
		HalfMoveClock = 0
		FullMoveCounter = 0
	}
}
