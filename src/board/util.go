package board

import (
	"strconv"
	"strings"
)

///////////////////////////////////////////////////////////////////
// General util
///////////////////////////////////////////////////////////////////

// setBit sets a bit on a bitboard
func setBit(bitboard *uint64, square int) {
	*bitboard |= (1 << uint64(square))
}

// getBit returns a bit on a bitboard
// mostly replace by indexMasks
func getBit(bitboard uint64, square int) uint64 {
	return bitboard & (1 << uint64(square))
}

// popBit pops a bit on a bitboard
func PopBit(bitboard *uint64, square int) {
	*bitboard ^= getBit(*bitboard, square)
}

// isBitOn returns true if a bit is on
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

// GetPieceType returns the piece type on a square
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

// GetPieceString returns the piece name on a square
func (cb *ChessBoard) GetPieceString(square int) string {
	for _, p := range AllPieceNames {
		bitboard := cb.GetPiecesBitboardString(p)
		if isBitOn(bitboard, square) {
			return p
		}
	}
	return "Empty"
}

// GetPiecesBitboardString returns the bitboard from a piece name
// for example: "WhitePawns" returns the bitboard of white pawns
// this is used to get the bitboard of a piece on a square
func (cb *ChessBoard) GetPiecesBitboardString(piece string) uint64 {
	pieceMap := map[string]uint64{
		"WhitePawns": cb.WhitePawns, "WhiteKnights": cb.WhiteKnights, "WhiteBishops": cb.WhiteBishops,
		"WhiteRooks": cb.WhiteRooks, "WhiteKing": cb.WhiteKing, "WhiteQueen": cb.WhiteQueen,

		"BlackPawns": cb.BlackPawns, "BlackKnights": cb.BlackKnights, "BlackBishops": cb.BlackBishops,
		"BlackRooks": cb.BlackRooks, "BlackKing": cb.BlackKing, "BlackQueen": cb.BlackQueen,
	}

	return pieceMap[piece]
}

// GetPieceInt returns the piece constant on a square
// for example: 1 is the bitboard for the white pawns
func (cb *ChessBoard) GetPieceInt(square int) int {
	for i := 0; i <= BlackKing; i++ {
		bitboard := cb.GetPiecesBitboardInt(i)
		if isBitOn(bitboard, square) {
			return i
		}
	}
	return Empty
}

// GetPiecesBitboardInt returns the bitboard from a piece constant
// for example: 1 is the bitboard for the white pawns
// constants are defined in constsvars.go
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
func (cb *ChessBoard) ParseFen(fen string) (err error) {
	*cb = ChessBoard{}
	// set default values
	SideToMove, CastleRights, Enpassant, HalfMoveClock, FullMoveCounter = 0, 0, -1, 0, 1
	fenRep := strings.Fields(fen)
	var file int
	var rank = 7

	// map each piece to its bitboard useful for updating the bitboards
	pieceMap := map[rune]*uint64{
		'P': &cb.WhitePawns, 'N': &cb.WhiteKnights, 'B': &cb.WhiteBishops, 'R': &cb.WhiteRooks, 'K': &cb.WhiteKing, 'Q': &cb.WhiteQueen,
		'p': &cb.BlackPawns, 'n': &cb.BlackKnights, 'b': &cb.BlackBishops, 'r': &cb.BlackRooks, 'k': &cb.BlackKing, 'q': &cb.BlackQueen,
	}

	// run through the fen string and update the bitboard for each piece
	for _, val := range fenRep[0] {
		switch val {
		// if the character is a slash, reset the file and decrement the rank
		case '/':
			file = 0
			rank--
		// if the character is a number, add that number to the file
		case '1', '2', '3', '4', '5', '6', '7', '8':
			file += int(val - '0')
		// if the character is a piece, set the bitboard for that piece
		default:
			setBit(pieceMap[val], (8*rank)+file)
			file++
		}
	}

	// set the white and black pieces bitboard
	cb.WhitePieces = cb.WhitePawns | cb.WhiteKnights | cb.WhiteBishops | cb.WhiteRooks | cb.WhiteQueen | cb.WhiteKing
	cb.BlackPieces = cb.BlackPawns | cb.BlackKnights | cb.BlackBishops | cb.BlackRooks | cb.BlackQueen | cb.BlackKing

	// set the side to move
	SideToMove = sideToMoveMap[fenRep[1]]

	// set the castle rights using CastleMap
	for _, val := range fenRep[2] {
		CastleRights += castleMap[val]
	}

	// set the enpassant square
	Enpassant = SquareToIndex[fenRep[3]]

	// set the halfmove clock and fullmove counter if they exist
	if len(fenRep) > 4 {
		HalfMoveClock, _ = strconv.Atoi(fenRep[4])
		FullMoveCounter, _ = strconv.Atoi(fenRep[5])
	} else {
		HalfMoveClock = 0
		FullMoveCounter = 0
	}

	// set the hash key
	HashKey = GenHash(*cb)

	return err
}
