package uci

import (
	"bufio"
	"os"

	"github.com/Droanox/ChessEngineAI/src/board"
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
		scan(cmd, &cb)
	}
}
