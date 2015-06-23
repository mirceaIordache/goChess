package quiesce

import (
	. "github.com/mirceaIordache/goChess/attack"
	. "github.com/mirceaIordache/goChess/common"
	. "github.com/mirceaIordache/goChess/distributed"
	. "github.com/mirceaIordache/goChess/evaluation"
	. "github.com/mirceaIordache/goChess/moveGenerator"
	. "github.com/mirceaIordache/goChess/sort"

	"strconv"
	"sync"
)

func LocalQuiesce(board ChessBoard, alpha, beta, depth int) int {
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s, alpha %d, beta %d, depth %d", ToEPD(board), alpha, beta, depth)

	if EvaluateDraw(board) == true || depth > MaxPlyDepth {
		ChessLogger.Info("Retruning a draw")
		ChessLogger.Debug("Result %d", SCORE_DRAW)
		return SCORE_DRAW
	}

	side := board.Side
	xside := 1 ^ side
	InCheck := SquareAttacked(board, board.KingPos[side], xside)
	best := Evaluate(alpha, beta, board)
	if best >= beta && InCheck == false {
		ChessLogger.Info("Found something above beta")
		ChessLogger.Debug("Result %d", best)
		return best
	}

	var list *MoveList

	if InCheck {
		ChessLogger.Info("In Check")
		list = GenerateCheckEscapes(board)
		if list == nil {
			ChessLogger.Info("No possible escapes")
			ChessLogger.Debug("Result %d", -Mate)
			return -Mate
		}
		if best >= beta {
			ChessLogger.Info("Found something above beta. In check")
			ChessLogger.Debug("Result %d", best)
			return best
		}

		SortMoves(board, list)
	} else {
		list = GenerateCaptures(board)
		if list == nil {
			ChessLogger.Info("No possible moves")
			ChessLogger.Debug("Result %d", best)
			return best
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

	done := make(chan ThreadedStatus, 100)
	var wg sync.WaitGroup
	counter := 0
	for list.Next != nil {
		p := PickBest(list)
		list = list.Next

		counter += 1
		ChessLogger.Debug("WG size: %d", counter)

		wg.Add(1)
		go func(p Move) {
			ChessLogger.Info("Started a new quiesce goroutine")
			if InCheck == false && SwapOff(board, p.Move) < delta {
				ChessLogger.Info("Quiesce Goroutine: Exiting, swap off below delta %d", delta)

				//				ChessLogger.Debug("Quiesce Goroutine: Result: %d", score)

				//				done <- ThreadedStatus{true, score}
				counter -= 1
				ChessLogger.Debug("Quiesce Goroutine: WG size: %d", counter)

				wg.Done()
				return
			}
			if p.Score == -Mate {
				ChessLogger.Info("Quiesce Goroutine: Exiting, move leads to a mate")
				counter -= 1
				ChessLogger.Debug("Quiesce Goroutine: WG size: %d", counter)

				wg.Done()
				return
			}
			newBoard := ApplyMove(board, p.Move, side)

			if SquareAttacked(newBoard, newBoard.KingPos[side], xside) {

				ChessLogger.Info("Quiesce Goroutine: Exiting, move leads to a check")
				counter -= 1
				ChessLogger.Debug("Quiesce Goroutine: WG size: %d", counter)

				wg.Done()
				return
			}
			score := -Quiesce(newBoard, -beta, -alpha, depth+1)
			if score > best {
				if score >= beta {
					ChessLogger.Info("Quiesce Goroutine: Exiting, found something above beta")
					ChessLogger.Debug("Quiesce Goroutine: Result: %d", score)

					done <- ThreadedStatus{true, score}
					counter -= 1
					ChessLogger.Debug("Quiesce Goroutine: WG size: %d", counter)

					wg.Done()
					return
				}
				if alpha > best {
					alpha = score
				}
			}
			done <- ThreadedStatus{true, score}

			counter -= 1
			ChessLogger.Info("No good result obtained, sending %d", score)
			ChessLogger.Debug("Quiesce Goroutine: WG size: %d", counter)

			wg.Done()
		}(p)
	}
	//	println("LocalQuiesce: ", ToEPD(board), alpha, beta, "  Result: ", sentinel.Value)
	wg.Wait()
	maxValue := -Mate
	ChessLogger.Info("Woke up from quiesce goroutines")

	sentinel := true
	for sentinel {
		select {
		case val := <-done:
			{
				if val.Status && val.Value > maxValue {
					maxValue = val.Value
				}
			}
		default:
			{
				sentinel = false
			}

		}
	}

	if best > maxValue {
		maxValue = best
	}

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result %d", maxValue)
	return maxValue
}

func Quiesce(board ChessBoard, alpha, beta, depth int) int {
	ChessLogger.Info("Going through proxy")
	res := LocalQuiesce(board, alpha, beta, depth)
	ChessLogger.Info("Proxy sent result")
	ChessLogger.Debug("Result %d", res)
	return res
	return RunRemote([]string{"Quiesce", "\"" + ToEPD(board) + "\"", strconv.Itoa(alpha), strconv.Itoa(beta), strconv.Itoa(depth)})

}
