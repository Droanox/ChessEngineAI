package main

import (
	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/graphics"
)

// command to run to create the .exe
// fyne package -os windows -icon ChessEngineAI.png
func main() {
	cb := board.ChessBoard{}
	cb.Init()
	graphics.Run(cb)
}
