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

	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			b := canvas.NewRectangle(color.RGBA{207, 167, 151, 1})
			if file%2 == rank%2 {
				b.FillColor = color.RGBA{150, 75, 45, 1}
			}

			//p := cb.Type(board.WhiteBishopsNum)
			piece := canvas.NewImageFromResource(pieceFromPNG(cb))
			piece.FillMode = canvas.ImageFillContain
			grid.Add(container.NewMax(b, piece))
		}
	}
	return grid
}
