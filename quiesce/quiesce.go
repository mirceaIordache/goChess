package quiesce

import (
	"fmt"
	. "github.com/mirceaIordache/goChess/attack"
	. "github.com/mirceaIordache/goChess/common"
	. "github.com/mirceaIordache/goChess/evaluation"
	. "github.com/mirceaIordache/goChess/moveGenerator"
	. "github.com/mirceaIordache/goChess/sort"
)

func Quiesce(board ChessBoard, alpha, beta, depth int) Move {

	iter := 1
	if EvaluateDraw(board) == true || depth > MaxPlyDepth {
		return Move{0, SCORE_DRAW}
	}

	side := board.Side
	xside := 1 ^ side
	InCheck := SquareAttacked(board, board.KingPos[side], xside)
	best := Evaluate(alpha, beta, board)
	if best >= beta && InCheck == false {

		return Move{0, best}
	}

	var list *MoveList

	if InCheck {
		list = GenerateCheckEscapes(board)
		FilterMoves(board, list, alpha, beta)
		if list == nil {
			return Move{0, -Mate}
		}
		if best >= beta {
			return Move{0, best}
		}
		SortMoves(board, list)
	} else {
		list = GenerateCaptures(board)
		if list == nil {
			return Move{0, best}
		}
		SortCaptures(board, list)
	}

	if best > alpha {
		alpha = best
	}
	delta := alpha - 150 - best
	if 0 > delta {
		delta = 0
	}

	var bestMove Move
	for list.Next != nil {
		verifyIntegrity(list, depth)
		p := PickBest(&list)
		if InCheck == false && SwapOff(board, p.Move) < delta {
			list = list.Next
			continue
		}

		if p.Score == -32767 {
			list = list.Next
			continue
		}
		newBoard := ApplyMove(board, p.Move, side)
		if SquareAttacked(newBoard, newBoard.KingPos[side], xside) {
			list = list.Next
			continue
		}

		score := -Quiesce(newBoard, -beta, -alpha, depth+1).Score
		if score > best {
			best = score
			bestMove = p
			if best >= beta {
				break
			}
			if alpha > best {
				alpha = best
			}
		}
		list = list.Next
		iter++
	}
	return bestMove
}

func verifyIntegrity(iter *MoveList, depth int) {
	for iter.Next != nil {
		//fmt.Println(iter.Value)

		if iter.Next == iter {
			fmt.Println("Depth: ", depth)
			panic("Circular ")
		}
		iter = iter.Next
	}
	//	fmt.Println("Integrity  OK")
	//	fmt.Println("---------------------")
}
