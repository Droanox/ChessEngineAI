//go:generate fyne bundle -o allPieces.go --pkg graphics pieces

package graphics

import (
	"fyne.io/fyne/v2"
	"github.com/Droanox/ChessEngineAI/src/board"
)

func pieceFromPNG(cb board.ChessBoard) fyne.Resource {
	switch {
	case cb.Type(board.WhitePawnsNum) == "WhitePawns":
		return resourceWhitePawnPng
	case cb.Type(board.WhiteKnightsNum) == "WhiteKnight":
		return resourceWhiteKnightPng
	case cb.Type(board.WhiteBishopsNum) == "WhiteBishop":
		return resourceWhiteBishopPng
	case cb.Type(board.WhiteRooksNum) == "WhiteRook":
		return resourceWhiteRookPng
	case cb.Type(board.WhiteQueensNum) == "WhiteQueen":
		return resourceWhiteQueenPng
	case cb.Type(board.WhiteKingNum) == "WhiteKing":
		return resourceWhiteKingPng

	case cb.Type(board.BlackPawnsNum) == "BlackPawns":
		return resourceBlackPawnPng
	case cb.Type(board.BlackKnightsNum) == "BlackKnight":
		return resourceBlackKnightPng
	case cb.Type(board.BlackBishopsNum) == "BlackBishop":
		return resourceBlackBishopPng
	case cb.Type(board.BlackRooksNum) == "BlackRook":
		return resourceBlackRookPng
	case cb.Type(board.BlackQueensNum) == "BlackQueen":
		return resourceBlackQueenPng
	case cb.Type(board.BlackKingNum) == "BlackKing":
		return resourceBlackKingPng
	}

	return nil
}
