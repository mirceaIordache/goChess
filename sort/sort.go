package sort

import (
	. "github.com/mirceaIordache/goChess/attack"
	. "github.com/mirceaIordache/goChess/common"
	. "github.com/mirceaIordache/goChess/evaluation"
	. "github.com/mirceaIordache/goChess/moveGenerator"
)

func SortMoves(board ChessBoard, list *MoveList) {
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s", ToEPD(board))
	side := board.Side
	xside := 1 ^ side
	enemyPawns := board.Board[xside][Pawn]
	iter := list
	cboard := GenerateCBoard(board)

	baseEval := Evaluate(-Mate, Mate, board)

	for iter.Next != nil {
		iter.Value.Score = baseEval
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

	ChessLogger.Info("Exiting")
}

func SortCaptures(board ChessBoard, list *MoveList) {
	/* Assign scores for captures only, no sorting */
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s", ToEPD(board))
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
	ChessLogger.Info("Exiting")
}

func PickBest(head *MoveList) Move {
	ChessLogger.Info("Entering")
	var prevInList *MoveList = nil
	var prevBest *MoveList = nil
	iter := head
	bestMove := iter
	prevInList = iter

	if iter.Next != nil {
		iter = iter.Next
	}

	for iter.Next != nil {
		if iter == iter.Next {
			ChessLogger.Error("List iter value: %d", iter.Value)
			ChessLogger.Panic("List is recursive!!!")
		}
		if iter.Value.Score > bestMove.Value.Score {
			bestMove = iter
			prevBest = prevInList
		}
		prevInList = iter
		iter = iter.Next
	}
	ChessLogger.Info("Exited the search loop")
	if bestMove != head {
		tmp := head.Next
		head.Next = bestMove.Next
		bestMove.Next = tmp
		prevBest.Next = head
		head = bestMove
	}

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result %d", bestMove.Value)
	return bestMove.Value

}
