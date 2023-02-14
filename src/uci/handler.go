package uci

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Droanox/ChessEngineAI/src/board"
	"github.com/Droanox/ChessEngineAI/src/search"
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
		cb.PrintChessBoard()
	case "position":
		handlePosition(commands, cb)
		cb.PrintChessBoard()
	case "go":
		handleGo(commands, cb)
	default:
		fmt.Println("invalid command")
	}
}

func handleUci() {
	fmt.Printf("id name ChessEngineAI\n")
	fmt.Printf("id author Leon Szabo\n")
	fmt.Printf("uciok\n")
}

func handleIsready() {
	fmt.Println("readyok")
}

func handleUcinewgame(cb *board.ChessBoard) {
	handlePosition("position startpos", cb)
}

func handlePosition(cmd string, cb *board.ChessBoard) {
	posCommands := strings.Fields(cmd)
	movesIndex := strings.Index(cmd, "moves")

	switch strings.ToLower(posCommands[1]) {
	case "startpos":
		cb.ParseFen(board.InitialPositionFen)
	case "fen":
		if movesIndex == -1 {
			cb.ParseFen(cmd[13:])
		} else {
			cb.ParseFen(cmd[13 : movesIndex-1])
		}
	default:
		fmt.Println("ERROR!")
		cb.ParseFen(board.InitialPositionFen)
	}

	if movesIndex != -1 {
		moveList := strings.Fields(cmd[movesIndex:])
		for _, move := range moveList {
			handleMakeMove(move, cb)
		}
	}
}

func handleGo(cmd string, cb *board.ChessBoard) {
	var depth int = 6

	goCommands := strings.Fields(cmd)

	if strings.ToLower(goCommands[1]) == "depth" {
		depth, _ = strconv.Atoi(goCommands[2])
	}

	search.Search(depth, cb)
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
					return true
				} else if ((flags & ^board.MoveCaptures) == board.MoveBishopPromotion) && move[4] == 'b' {
					cb.MakeMove(moveList[i])
					return true
				} else if ((flags & ^board.MoveCaptures) == board.MoveRookPromotion) && move[4] == 'r' {
					cb.MakeMove(moveList[i])
					return true
				} else if ((flags & ^board.MoveCaptures) == board.MoveQueenPromotion) && move[4] == 'q' {
					cb.MakeMove(moveList[i])
					return true
				}
				continue
			}
			cb.MakeMove(moveList[i])
			board.Ply = -1
			return true
		}
	}
	return false
}
