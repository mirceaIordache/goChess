package main

import (
	. "github.com/mirceaIordache/goChess/common"
	. "github.com/mirceaIordache/goChess/evaluation"
	. "github.com/mirceaIordache/goChess/quiesce"
	. "github.com/mirceaIordache/goChess/search"

	"github.com/op/go-logging"

	"fmt"
	"os"
	"strconv"
)

var version string

func main() {
	init_Logger(false, true)

	ChessLogger.Debug("goChess version %s started", version)

	switch len(os.Args) {
	//initial invocation
	case 3:
		main_app(os.Stdout)
	//Remote Quiescence
	case 7:
		handle_quiesce()

	//Remote Search
	case 8:
		handle_search()

	//Stuff
	default:
		fmt.Println("Usage: goChess circuit://... EPD")
	}
}

func main_app(stream *os.File) {
	ChessLogger.Info("Started as main application")
	ChessLogger.Debug("os.Args is %q", os.Args[1:len(os.Args)])
	var rootAlpha, rootBeta, iterationDepth int
	var bestMove = Move{-Mate, -Mate}
	ChessLogger.Debug("BestMove set to %d", bestMove)
	iterationDepth = 1
	boardStr := os.Args[2]

	board := FromEPD(boardStr)

	score := Evaluate(-Mate, Mate, board)

	if score > Mate-255 {
		rootAlpha = score - 1
		rootBeta = Mate
	} else if score < -Mate+255 {
		rootAlpha = -Mate
		rootBeta = score + 1
	} else {
		if score-75 > Mate {
			rootAlpha = score - 75
		} else {
			rootAlpha = Mate
		}
		if score+75 < -Mate {
			rootBeta = score + 75
		} else {
			rootBeta = -Mate
		}
	}
	for iterationDepth = 1; iterationDepth <= 10; iterationDepth++ {
		ChessLogger.Debug("Iteration depth %d", iterationDepth)
		move := SearchRoot(board, iterationDepth, rootAlpha, rootBeta)
		if move.Score >= rootBeta && move.Score < Mate {
			SetBest(-Mate + 1)
			rootAlpha = rootBeta
			rootBeta = Mate
			move = SearchRoot(board, iterationDepth, rootAlpha, rootBeta)
		} else if move.Score <= rootAlpha {
			SetBest(-Mate + 1)
			rootBeta = rootAlpha
			rootAlpha = -Mate
			move = SearchRoot(board, iterationDepth, rootAlpha, rootBeta)
		}

		if move.Score > bestMove.Score {
			SetBest(move.Score)
			bestMove = move
			ChessLogger.Debug("BestMove set to %d", bestMove)
		}
		ChessLogger.Debug("Best Move found so far is %s-%s with a score of %d", Algebraic[MoveFrom(bestMove.Move)], Algebraic[MoveTo(bestMove.Move)], bestMove.Score)
	}

	fmt.Fprintf(stream, "Best move is: %s-%s with a score of %d \n", Algebraic[MoveFrom(bestMove.Move)], Algebraic[MoveTo(bestMove.Move)], bestMove.Score)
	ChessLogger.Info("Done searching for the move. Results printed out")
}

func handle_search() {
	ChessLogger.Info("Started as a remote searcher")
	ChessLogger.Debug("os.Args is %q", os.Args[2:len(os.Args)])
	//func LocalSearch(board ChessBoard, depth, alpha, beta, nodeType int) int
	board := FromEPD(os.Args[3])
	depth, _ := strconv.Atoi(os.Args[4])
	alpha, _ := strconv.Atoi(os.Args[5])
	beta, _ := strconv.Atoi(os.Args[6])
	nodeType, _ := strconv.Atoi(os.Args[7])

	move := LocalSearch(board, depth, alpha, beta, nodeType)

	fmt.Println(strconv.Itoa(move))
	ChessLogger.Debug("Result is %s-%s", Algebraic[MoveFrom(move)], Algebraic[MoveTo(move)])
	ChessLogger.Info("Done searching, sent result to parent")
}

func handle_quiesce() {
	//func Quiesce(board ChessBoard, alpha, beta int) int
	ChessLogger.Info("Started as a remote quiescer")
	ChessLogger.Debug("os.Args is %q", os.Args[2:len(os.Args)])
	board := FromEPD(os.Args[3])
	alpha, _ := strconv.Atoi(os.Args[4])
	beta, _ := strconv.Atoi(os.Args[5])
	depth, _ := strconv.Atoi(os.Args[6])

	move := LocalQuiesce(board, alpha, beta, depth)

	fmt.Println(strconv.Itoa(move))
	ChessLogger.Debug("Result is %s-%s", Algebraic[MoveFrom(move)], Algebraic[MoveTo(move)])
	ChessLogger.Info("Done quiescing, sent result to parent")
}

func init_Logger(colour, enabled bool) {
	ChessLogger = logging.MustGetLogger("goChess")

	colourSet := ""
	colourReset := ""
	if colour {
		colourSet = "%{color}"
		colourReset = "%{color:reset}"
	}
	backend, _ := os.Open(os.DevNull)
	if enabled {
		backend = os.Stderr
	}
	format := logging.MustStringFormatter(colourSet + "%{time:15:04:05.000} %{id:d} %{level:s} -> %{shortfunc} :" + colourReset + " %{message}")

	generalBackend := logging.NewLogBackend(backend, "", 0)
	generalFormatter := logging.NewBackendFormatter(generalBackend, format)

	errorLeveled := logging.AddModuleLevel(generalBackend)
	errorLeveled.SetLevel(logging.ERROR, "")

	logging.SetBackend(errorLeveled, generalFormatter)
}
