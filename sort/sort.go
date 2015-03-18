package sort

import (
	"fmt"
	//	"runtime"

	. "github.com/mirceaIordache/goChess/attack"
	. "github.com/mirceaIordache/goChess/common"
	. "github.com/mirceaIordache/goChess/moveGenerator"
)

func SortMoves(board ChessBoard, list *MoveList) {
	side := board.Side
	xside := 1 ^ side
	enemyPawns := board.Board[xside][Pawn]
	iter := list
	cboard := GenerateCBoard(board)

	for iter.Next != nil {
		from := MoveFrom(iter.Value.Move)
		to := MoveTo(iter.Value.Move)

		if cboard[to] != Empty || iter.Value.Move&PiecePromotion != 0 {
			toValue := Values[cboard[to]] + Values[PromotePiece(iter.Value.Move)]
			iter.Value.Score += int(toValue) + int(ValueKing) - int(Values[cboard[from]])
		}

		if cboard[from] == Pawn {
			if enemyPawns&PassedPawnMask[side][to] == NullBitBoard {
				iter.Value.Score += 50
			}
		}

		iter = iter.Next
	}
}

func SortCaptures(board ChessBoard, list *MoveList) {
	/* Assign scores for captures only, no sorting */
	iter := list
	for iter.Next != nil {
		cboard := GenerateCBoard(board)
		from := Values[cboard[MoveFrom(iter.Value.Move)]]
		to := Values[cboard[MoveTo(iter.Value.Move)]]
		if from < to {
			iter.Value.Score = int(to) - int(from)
		} else {
			temp := SwapOff(board, iter.Value.Move)
			if temp < 0 {
				temp = -Mate
			}
			iter.Value.Score = temp
		}
		iter = iter.Next
	}
}

func PickBest(head **MoveList) Move {
	var prevInList *MoveList = nil
	var prevBest *MoveList = nil
	iter := *head
	best := iter.Value.Score
	bestMove := iter

	for iter.Next != nil {

		if iter.Next == iter {
			fmt.Println("Proof: ", iter.Next.Value)
			panic("Circular")
		}
		iter = iter.Next
	}
	iter = *head
	prevInList = iter
	iter = iter.Next
	for iter.Next != nil {
		if iter.Value.Score > best {
			bestMove = iter
			best = iter.Value.Score
			prevBest = prevInList
		}
		prevInList = iter
		iter = iter.Next
	}

	if bestMove != *head && bestMove != (*head).Next {
		tmp := (*head).Next
		(*head).Next = bestMove.Next
		bestMove.Next = tmp
		prevBest.Next = *head
		*head = bestMove
	} else if bestMove == (*head).Next {
		(*head).Next = bestMove.Next
		bestMove.Next = *head
		*head = bestMove
	}

	return bestMove.Value

}
