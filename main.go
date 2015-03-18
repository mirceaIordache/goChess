package main

import (
	"fmt"

	. "github.com/mirceaIordache/goChess/common"
	. "github.com/mirceaIordache/goChess/evaluation"
	. "github.com/mirceaIordache/goChess/search"
)

func main() {
	var rootAlpha, rootBeta, iterationDepth int
	var bestMove Move
	iterationDepth = 0
	//Tactical: 8.0
	//Lever: 0.0
	// Rating < 1500
	boardStr := "r2qnrnk/p2b2b1/1p1p2pp/2pPpp2/1PP1P3/PRNBB3/3QNPPP/5RK1 w - -"

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
	iterationDepth++
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

	SetBest(move.Score)
	bestMove = move

	fmt.Println("Best move is:", bestMove.Move, "(", Algebraic[MoveFrom(bestMove.Move)], ",", Algebraic[MoveTo(bestMove.Move)], ") with a score of", bestMove.Score)
}

/*
import (
	"fmt"
	. "github.com/mirceaIordache/goChess/common"
	. "github.com/mirceaIordache/goChess/moveGenerator"
	. "github.com/mirceaIordache/goChess/sort"
)

func main() {
	boardStr := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - -"
	board := FromEPD(boardStr)

	list := GenerateMoves(board)
	FilterMoves(board, list, Mate, -Mate)

	for list.Next != nil {
		PickBest(&list)
		list = list.Next
	}

	fmt.Println("Done")
	//	p.Score++
}
*/
