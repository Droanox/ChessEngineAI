package uci

import (
	"bufio"
	"os"

	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/engine"
	"github.com/Droanox/ChessEngineAI/src/eval"
)

func Run() {
	cb := board.ChessBoard{}
	cb.Init()
	eval.Init()

	// Make a channel to receive commands from the user
	cmdCh := make(chan string)
	// Make a bool to control the scanning of the user input
	var okayToScan bool = true
	// Concurrently scan the user input
	go scanLine(cmdCh, &okayToScan)

	for {
		// Wait for a command to be sent to the channel
		cmd := <-cmdCh
		// Set the bool to false so that the user input is not scanned
		okayToScan = false
		// Scan the command
		scan(cmd, &cb)
		// Set the bool to true so that the user input can be scanned again
		okayToScan = true
	}
}

func scanLine(cmdCh chan string, okayToScan *bool) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		switch cmd := scanner.Text(); cmd {
		case "quit":
			return
		case "stop":
			engine.IsStopped = true
		default:
			// if the engine is not searching, we can scan the command
			if (*okayToScan && cmd != "") || cmd == "isready" {
				cmdCh <- cmd
			}
		}
	}
}
