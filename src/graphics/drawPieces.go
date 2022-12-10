//go:generate fyne bundle -o allPieces.go --pkg graphics pieces

package graphics

import "fyne.io/fyne/v2"

func pieceFromPNG() fyne.Resource {
	return resourceBlackBishopPng
}
