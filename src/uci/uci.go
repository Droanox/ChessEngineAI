package uci

import (
	"bufio"
	"os"

	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/engine"
)

func Run() {
	scanner := bufio.NewScanner(os.Stdin)
	cb := board.ChessBoard{}
	cb.Init()

	for scanner.Scan() {
		var cmd = scanner.Text()
		if cmd == "quit" {
			break
		}
		if cmd == "stop" {
			engine.IsStopped = true
			continue
		}
		go scan(cmd, &cb)
	}
}
