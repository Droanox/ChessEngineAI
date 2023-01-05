//go:generate fyne bundle -o allPieces.go --pkg graphics pieces

package graphics

import (
	"fyne.io/fyne/v2"
)

func pieceFromPNG(piece string) fyne.Resource {
	switch piece {
	case "WhitePawns":
		return resourceWhitePawnPng
	case "WhiteKnights":
		return resourceWhiteKnightPng
	case "WhiteBishops":
		return resourceWhiteBishopPng
	case "WhiteRooks":
		return resourceWhiteRookPng
	case "WhiteQueen":
		return resourceWhiteQueenPng
	case "WhiteKing":
		return resourceWhiteKingPng

	case "BlackPawns":
		return resourceBlackPawnPng
	case "BlackKnights":
		return resourceBlackKnightPng
	case "BlackBishops":
		return resourceBlackBishopPng
	case "BlackRooks":
		return resourceBlackRookPng
	case "BlackQueen":
		return resourceBlackQueenPng
	case "BlackKing":
		return resourceBlackKingPng
	}

	return nil
}
