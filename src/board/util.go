package board

import (
	"reflect"
	"strings"
	"unicode"
)

func setBit(bitboard *uint64, square int) {
	*bitboard |= (1 << uint64(square))
}

func GetBit(bitboard uint64, square int) bool {
	return bitboard == (bitboard | (1 << uint64(square)))
}

func (cb *ChessBoard) GetPiece(square int) string {
	for _, p := range AllPieceNames {
		bitboard := cb.GetPiecesBitboard(p)
		if GetBit(bitboard, square) {
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

func (cb ChessBoard) Type(num int) string {
	b := reflect.TypeOf(cb)
	return b.Field(num).Name
}

// "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNN w KQkq - 0 1"
func (cb *ChessBoard) parseFen(fen string) {
	fenRep := strings.Split(fen, " ")
	file := 0
	rank := 7

	for _, val := range fenRep[0] {
		if val == '/' {
			file = 0
			rank--
		} else {
			if unicode.IsDigit(rune(val)) {
				file += (int(val) - '0')
			} else {
				conv := (8 * rank) + file
				switch val {
				// Black pieces
				case 114:
					setBit(&cb.BlackRooks, conv)
				case 110:
					setBit(&cb.BlackKnights, conv)
				case 98:
					setBit(&cb.BlackBishops, conv)
				case 113:
					setBit(&cb.BlackQueen, conv)
				case 107:
					setBit(&cb.BlackKing, conv)
				case 112:
					setBit(&cb.BlackPawns, conv)
				// White pieces
				case 82:
					setBit(&cb.WhiteRooks, conv)
				case 78:
					setBit(&cb.WhiteKnights, conv)
				case 66:
					setBit(&cb.WhiteBishops, conv)
				case 81:
					setBit(&cb.WhiteQueen, conv)
				case 75:
					setBit(&cb.WhiteKing, conv)
				case 80:
					setBit(&cb.WhitePawns, conv)
				}
				file++
			}
		}
	}
}
