package uci

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/engine"
)

func scan(commands string, cb *board.ChessBoard) {
	cmd := strings.Fields(commands)[0]

	switch cmd {
	case "uci":
		handleUci()
	case "isready":
		handleIsready()
	case "ucinewgame":
		handleUcinewgame(cb)
	case "position":
		isReady = false
		handlePosition(commands, cb)
		isReady = true
	case "go":
		handleGo(commands, cb)
	default:
		return
	}
}

func handleUci() {
	fmt.Printf("id name ChessEngineAI\n")
	fmt.Printf("id author Leon Szabo & Flávia Alves\n")
	fmt.Printf("uciok\n")
}

func handleIsready() {
	if isReady {
		fmt.Println("readyok")
	}
}

func handleUcinewgame(cb *board.ChessBoard) {
	*cb = board.ChessBoard{}

	engine.ClearTT()
}

func handlePosition(cmd string, cb *board.ChessBoard) (err error) {
	posCommands := strings.Fields(cmd)
	movesIndex := strings.Index(cmd, "moves")

	switch strings.ToLower(posCommands[1]) {
	case "startpos":
		cb.ParseFen(board.InitialPositionFen)
	case "fen":
		if movesIndex == -1 {
			err = cb.ParseFen(cmd[13:])
		} else {
			err = cb.ParseFen(cmd[13 : movesIndex-1])
		}
	default:
		fmt.Println("Position not found, using default position")
		cb.ParseFen(board.InitialPositionFen)
	}

	if movesIndex != -1 {
		moveList := strings.Fields(cmd[movesIndex:])
		for _, move := range moveList {
			handleMakeMove(move, cb)
		}
	}

	hasParsed = true

	// engine.ClearTT()

	return err
}

func handleGo(cmd string, cb *board.ChessBoard) (err error) {
	if !hasParsed {
		handlePosition("position startpos", cb)
	}
	var movesToGo, depth, stopTime int = 1, -1, 0
	var timeLeft time.Duration

	goCommands := strings.Fields(cmd)

	if len(goCommands) > 1 {
		for i, command := range goCommands {
			switch command {
			case "wtime":
				if board.SideToMove == board.White {
					timeLeft, err = time.ParseDuration(goCommands[i+1] + "ms")
				}
			case "btime":
				if board.SideToMove == board.Black {
					timeLeft, err = time.ParseDuration(goCommands[i+1] + "ms")
				}
			case "movestogo":
				movesToGo, err = strconv.Atoi(goCommands[i+1])
			case "depth":
				depth, err = strconv.Atoi(goCommands[i+1])
			case "movetime":
				stopTime, err = strconv.Atoi(goCommands[i+1])
				engine.TimeControl = true
				engine.StopTime = int64(stopTime)
			}
		}
	}

	if depth == -1 {
		depth = board.MaxPly
	}

	if stopTime == 0 && timeLeft != 0 {
		engine.TimeControl = true

		engine.StopTime = timeLeft.Milliseconds()/int64(movesToGo) - ((timeLeft.Milliseconds() / int64(movesToGo)) / 10)
	}
	engine.Search(depth, cb)

	return err
}

func handleMakeMove(move string, cb *board.ChessBoard) bool {
	var start int = board.SquareToIndex[move[0:2]]
	var end int = board.SquareToIndex[move[2:4]]
	var promotionMask int

	var moveList = []board.Move{}
	cb.GenerateMoves(&moveList)

	for i := 0; i < len(moveList); i++ {
		if (moveList[i].GetMoveStart() == start) && (moveList[i].GetMoveEnd() == end) {
			promotionMask = 0b1000 & moveList[i].GetMoveFlags()
			flags := moveList[i].GetMoveFlags()

			if promotionMask != 0 {
				if ((flags & ^board.MoveCaptures) == board.MoveKnightPromotion) && move[4] == 'n' {
					cb.MakeMove(moveList[i])
					board.RepetitionTableIndexOffset++
					board.Ply--
					return true
				} else if ((flags & ^board.MoveCaptures) == board.MoveBishopPromotion) && move[4] == 'b' {
					cb.MakeMove(moveList[i])
					board.RepetitionTableIndexOffset++
					board.Ply--
					return true
				} else if ((flags & ^board.MoveCaptures) == board.MoveRookPromotion) && move[4] == 'r' {
					cb.MakeMove(moveList[i])
					board.RepetitionTableIndexOffset++
					board.Ply--
					return true
				} else if ((flags & ^board.MoveCaptures) == board.MoveQueenPromotion) && move[4] == 'q' {
					cb.MakeMove(moveList[i])
					board.RepetitionTableIndexOffset++
					board.Ply--
					return true
				}
				continue
			} else {
				cb.MakeMove(moveList[i])
				board.RepetitionTableIndexOffset++
				board.Ply--
				return true
			}
		}
	}
	return false
}
