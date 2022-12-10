package graphics

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func Run() {
	app := app.New()
	w := app.NewWindow("Chess")

	grid := createBoard()
	w.SetContent(grid)
	w.Resize(fyne.NewSize(600, 600))

	w.ShowAndRun()
}

func createBoard() *fyne.Container {
	grid := container.NewGridWithColumns(8)

	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			board := canvas.NewRectangle(color.RGBA{207, 167, 151, 1})
			if file%2 == rank%2 {
				board.FillColor = color.RGBA{150, 75, 45, 1}
			}
			piece := canvas.NewImageFromResource(pieceFromPNG())
			piece.FillMode = canvas.ImageFillContain
			grid.Add(container.NewMax(board, piece))
		}
	}
	return grid
}
