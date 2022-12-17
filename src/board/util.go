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
	return bitboard == (1 << uint64(square))
}

func (cb ChessBoard) GetPiece(square int) string {
	/*
		for i := 0; i < 12; i++{

		}
		switch cb {
		case cb.BlackPawns:

		}
	*/
	return "Empty"
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
				file += int(val)
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
					setBit(&cb.BlackQueens, conv)
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
					setBit(&cb.WhiteQueens, conv)
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
