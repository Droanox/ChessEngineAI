package graphics

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/Droanox/ChessEngineAI/src/board"
)

func Run(cb board.ChessBoard) {
	app := app.New()
	w := app.NewWindow("Chess")

	grid := createBoard(cb)
	w.SetContent(grid)
	w.Resize(fyne.NewSize(600, 600))

	w.ShowAndRun()
}

func createBoard(cb board.ChessBoard) *fyne.Container {
	grid := container.NewGridWithColumns(8)

	for rank := 8; rank >= 1; rank-- {
		for file := 1; file <= 8; file++ {
			b := canvas.NewRectangle(color.RGBA{207, 167, 151, 255})
			if file%2 == rank%2 {
				b.FillColor = color.RGBA{150, 75, 45, 255}
			}
			piece := canvas.NewImageFromResource(pieceFromPNG(cb.GetPiece((rank-1)*8 + (file - 1))))
			piece.FillMode = canvas.ImageFillContain
			grid.Add(container.NewMax(b, piece))
		}
	}
	return grid
}
