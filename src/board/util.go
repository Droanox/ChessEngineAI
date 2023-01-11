package board

import (
	"strconv"
	"strings"
	"unicode"
)

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

// The "Brian Kernighan's way" of counting bits on a bitboard,
// implementation idea from chess programming wiki
func BitCount(bitboard uint64) int {
	count := 0
	for bitboard != EmptyBoard {
		count++
		bitboard &= bitboard - 1
	}
	return count
}

func BitScanForward(bitboard uint64) int {
	const debruijn64 uint64 = 0x03f79d71b4cb0a89
	if bitboard != 0 {
		return index64[((bitboard^(bitboard-1))*debruijn64)>>58]
	}
	return -1
}

func BitScanReverse(bitboard uint64) int {
	const debruijn64 uint64 = 0x03f79d71b4cb0a89
	bitboard |= bitboard >> 1
	bitboard |= bitboard >> 2
	bitboard |= bitboard >> 4
	bitboard |= bitboard >> 8
	bitboard |= bitboard >> 16
	bitboard |= bitboard >> 32
	if bitboard != 0 {
		return index64[(bitboard*debruijn64)>>58]
	}
	return -1
}

func (cb *ChessBoard) GetPiece(square int) string {
	for _, p := range AllPieceNames {
		bitboard := cb.GetPiecesBitboard(p)
		if isBitOn(bitboard, square) {
			return p
		}
	}
	return "Empty"
}

func (cb *ChessBoard) GetPiecesBitboard(p string) uint64 {
	switch p {
	case "WhitePawns":
		return cb.WhitePawns
	case "WhiteKnights":
		return cb.WhiteKnights
	case "WhiteBishops":
		return cb.WhiteBishops
	case "WhiteRooks":
		return cb.WhiteRooks
	case "WhiteQueen":
		return cb.WhiteQueen
	case "WhiteKing":
		return cb.WhiteKing

	case "BlackPawns":
		return cb.BlackPawns
	case "BlackKnights":
		return cb.BlackKnights
	case "BlackBishops":
		return cb.BlackBishops
	case "BlackRooks":
		return cb.BlackRooks
	case "BlackQueen":
		return cb.BlackQueen
	case "BlackKing":
		return cb.BlackKing
	}
	return EmptyBoard
}

/* Was used but made redundant, code is kept if it's ever needed again
func (cb ChessBoard) Type(num int) string {
	b := reflect.TypeOf(cb)
	return b.Field(num).Name
}
*/

// Parses a fen string for example: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNN w KQkq - 0 1"
// onto the chessboard and maps every pieces bitboard to the relevant pieces
func (cb *ChessBoard) parseFen(fen string) {
	fenRep := strings.Split(fen, " ")
	file := 0
	rank := 7
	square := 0

	for _, val := range fenRep[0] {
		if val == '/' {
			file = 0
			rank--
		} else {
			if unicode.IsDigit(rune(val)) {
				file += (int(val) - '0')
			} else {
				square = (8 * rank) + file
				switch val {
				// Black pieces
				case 114:
					setBit(&cb.BlackRooks, square)
				case 110:
					setBit(&cb.BlackKnights, square)
				case 98:
					setBit(&cb.BlackBishops, square)
				case 113:
					setBit(&cb.BlackQueen, square)
				case 107:
					setBit(&cb.BlackKing, square)
				case 112:
					setBit(&cb.BlackPawns, square)
				// White pieces
				case 82:
					setBit(&cb.WhiteRooks, square)
				case 78:
					setBit(&cb.WhiteKnights, square)
				case 66:
					setBit(&cb.WhiteBishops, square)
				case 81:
					setBit(&cb.WhiteQueen, square)
				case 75:
					setBit(&cb.WhiteKing, square)
				case 80:
					setBit(&cb.WhitePawns, square)
				}
				file++
			}
		}
	}
	cb.WhitePieces = cb.WhiteRooks | cb.WhiteKnights | cb.WhiteBishops | cb.WhiteQueen | cb.WhiteKing | cb.WhitePawns
	cb.BlackPieces = cb.BlackRooks | cb.BlackKnights | cb.BlackBishops | cb.BlackQueen | cb.BlackKing | cb.BlackPawns

	if fenRep[1] == "w" {
		SideToMove = White
	} else {
		SideToMove = Black
	}

	for _, val := range fenRep[2] {
		switch val {
		case 'K':
			CastleRights += WhiteKingSide
		case 'Q':
			CastleRights += WhiteQueenSide
		case 'k':
			CastleRights += BlackKingSide
		case 'q':
			CastleRights += BlackQueenSide
		}
	}

	if fenRep[3] == "-" {
		Enpassant = -1
	} else {
		Enpassant = SquareToInt[fenRep[3]]
	}

	if len(fenRep) > 4 {
		HalfMoveClock, _ = strconv.Atoi(fenRep[4])
		FullMoveCounter, _ = strconv.Atoi(fenRep[5])
	} else {
		HalfMoveClock = 0
		FullMoveCounter = 0
	}
}
