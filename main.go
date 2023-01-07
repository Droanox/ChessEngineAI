package main

import (
	"github.com/Droanox/ChessEngineAI/src/board"
)

// command to run to create the .exe
// fyne package -os windows -icon ChessEngineAI.png
func main() {
	cb := board.ChessBoard{}
	cb.Init()
	cb.Test()
	//graphics.Run(cb)
}
