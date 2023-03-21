package uci

import (
	"bufio"
	"os"
	"strings"

	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/engine"
	"github.com/Droanox/ChessEngineAI/src/eval"
)

// validCommands is a map of valid commands
var validCommands = map[string]bool{"uci": true, "isready": true, "ucinewgame": true, "position": true, "ponderhit": true, "go": true}

// isReady is a bool to check if the engine is ready
var isReady bool = true

// hasParsed is a bool to check if the engine has parsed the position
// if not, it will use the initial position
var hasParsed bool = false

// chMax is the maximum size of the channel
var chMax int = 10

func Run() {
	cb := board.ChessBoard{}
	cb.Init()
	eval.Init()

	// Make a channel to receive commands from the user
	cmdCh := make(chan string, chMax)
	// Make a bool to control the scanning of the user input
	var okayToScan bool = true
	// Concurrently scan the user input
	go scanLine(cmdCh, &okayToScan)

	for {
		// Wait for a command to be sent to the channel
		cmd := <-cmdCh
		// If the command is "quit", return
		if cmd == "quit" {
			return
		}
		// Scan the command
		scan(cmd, &cb)
	}
}

func scanLine(cmdCh chan string, okayToScan *bool) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		switch cmd := scanner.Text(); cmd {
		case "quit":
			// If the engine is searching, we need to stop it first
			engine.IsStopped = true
			// Empty the channel
			for len(cmdCh) > 0 {
				<-cmdCh
			}
			// Send the quit command to the channel
			cmdCh <- "quit"
			return
		case "stop":
			engine.IsStopped = true
		default:
			// if the engine is not searching, we can scan the command
			if cmd != "" && validCommands[strings.Fields(cmd)[0]] {
				if len(cmdCh) == chMax {
					<-cmdCh
				} else {
					cmdCh <- cmd
				}
			}
		}
	}
}
