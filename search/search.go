package search

import (
	. "github.com/mirceaIordache/goChess/attack"
	. "github.com/mirceaIordache/goChess/common"
	. "github.com/mirceaIordache/goChess/distributed"
	. "github.com/mirceaIordache/goChess/evaluation"
	. "github.com/mirceaIordache/goChess/moveGenerator"
	. "github.com/mirceaIordache/goChess/quiesce"
	. "github.com/mirceaIordache/goChess/sort"

	"strconv"
	"sync"
)

var overallBest = -Mate

func SetBest(best int) {
	overallBest = best
}

func SearchRoot(board ChessBoard, iterationDepth, alpha, beta int) Move {

	var score int
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s, depth %d, alpha %d, beta %d", ToEPD(board), iterationDepth, alpha, beta)
	side := board.Side
	best := -Mate
	firstMove := true
	nodeType := PV
	bestMove := Move{}
	list := GenerateMoves(board)
	FilterMoves(board, list)
	SortMoves(board, list)
	iter := list

	for iter.Next != nil {
		p := PickBest(iter)
		newBoard := ApplyMove(board, p.Move, side)
		if p == list.Value {
			score = -Search(newBoard, iterationDepth-1, -beta, -alpha, nodeType, firstMove)
			if beta == Mate && score <= alpha {
				alpha = -Mate
				score = -Search(board, iterationDepth-1, -beta, -alpha, nodeType, firstMove)
			}
		} else {
			nodeType = CUT
			if best > alpha {
				alpha = best
			}
			score = -Search(board, iterationDepth-1, -alpha-1, -alpha, nodeType, firstMove)
			if score > best {
				if alpha < score && score < beta {
					nodeType = PV
					score = -Search(board, iterationDepth-1, -beta, -score, nodeType, firstMove)
				}
			}
		}
		iter.Value.Score = score
		if score > best {
			best = score
			bestMove = p
			if best > alpha {
				overallBest = best
				if best >= beta {
					break
				}
			}
		}
		if Mate+1 == best+1 {
			ChessLogger.Info("Exiting")
			ChessLogger.Debug("70: Result %d", bestMove)
			return bestMove
		}
		iter = iter.Next
		firstMove = false
	}

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("77: Result %d", bestMove)
	return bestMove
}

func LocalSearch(board ChessBoard, depth, alpha, beta, nodeType int) int {
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s, depth %d, alpha %d, beta %d, nodetype %d", ToEPD(board), depth, alpha, beta, nodeType)

	if EvaluateDraw(board) == true {
		ChessLogger.Info("Found a draw, exiting")
		ChessLogger.Debug("87: Result %d", SCORE_DRAW)
		return SCORE_DRAW
	}

	list := GenerateMoves(board)
	FilterMoves(board, list)
	SortMoves(board, list)

	var score int
	best := -Mate
	side := board.Side
	xside := 1 ^ side
	inCheck := SquareAttacked(board, board.KingPos[side], xside)
	if inCheck == true {
		list = GenerateCheckEscapes(board)
		if list.Next == nil {
			ChessLogger.Info("Mate. Exiting")
			ChessLogger.Debug("103: Result %d", Mate)
			return Mate
		}
	}

	if overallBest >= Mate {
		res := CalculateMaterial(board)
		ChessLogger.Info("Huge advantage so far")
		ChessLogger.Debug("111: Result %d", res)
		return res
	}

	if depth <= 0 {
		ChessLogger.Info("Reached search depth. Quiescing")

		res := Quiesce(board, alpha, beta, 1)
		ChessLogger.Debug("119: Result %d", res)
		return res
	}

	if depth > 1 && nodeType != PV && inCheck == false && CalculateMaterial(board)+ValuePawn > beta && beta > -Mate && board.PawnMaterial[side] > ValueBishop {
		newBoard := ApplyNullMove(board, 0, side)
		nullScore := -Search(newBoard, depth-3, -beta, -beta+1, nodeType, true)
		if nullScore >= beta {
			ChessLogger.Info("Null move is best")
			ChessLogger.Debug("128: Result %d", nullScore)
			return nullScore
		}
		if depth-3 >= 1 && CalculateMaterial(board) > beta && nullScore <= -Mate+256 {
			depth++
		}
	}

	if inCheck {
		SortMoves(board, list)
	}

	firstMove := true
	saveNodeType := nodeType

	p := PickBest(list)
	newBoard := ApplyMove(board, p.Move, side)
	score = -Search(newBoard, depth-1, -beta, -alpha, nodeType, firstMove)
	firstMove = false

	if score > best {
		best = score
		if best >= beta {
			return best
		}
	}

	if Mate+1 == best {
		return best
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
			ChessLogger.Info("New goroutine started")
			newBoard := ApplyMove(board, p.Move, side)
			if list.Next == nil {
				ChessLogger.Info("Search Goroutine: no more moves, exiting")
				ChessLogger.Debug("Search Goroutine: Result %d", best)
				done <- ThreadedStatus{true, best}
				counter -= 1
				ChessLogger.Debug("Quiesce Goroutine: WG size: %d", counter)

				wg.Done()
				return
			}

			if SquareAttacked(newBoard, newBoard.KingPos[side], xside) {
				ChessLogger.Info("Search Goroutine: check. Aborting")
				counter -= 1
				ChessLogger.Debug("Quiesce Goroutine: WG size: %d", counter)

				wg.Done()
				return
			}

			if nodeType == PV {
				nodeType = CUT
			}

			if best > alpha {
				alpha = best
			}

			score = -Search(newBoard, depth-1, -alpha-1, -alpha, nodeType, firstMove)
			if score > best {
				if saveNodeType == PV {
					nodeType = PV
				}
				if alpha < score && score < beta {
					score = -Search(newBoard, depth-1, -beta, -score, nodeType, firstMove)
				}
				if nodeType == PV && score <= alpha {
					score = -Search(newBoard, depth-1, -alpha, -Mate, nodeType, firstMove)
				}
			}

			if score > best {
				best = score
				if best >= beta {
					ChessLogger.Info("Search Goroutine: Found new best value")
					ChessLogger.Debug("Search Goroutine: Result %d", score)
					counter -= 1
					ChessLogger.Debug("Quiesce Goroutine: WG size: %d", counter)

					wg.Done()
					done <- ThreadedStatus{true, score}
				}
			}

			if Mate+1 == best {
				ChessLogger.Info("Search Goroutine: Apparently doing a mate")
				ChessLogger.Debug("Search Goroutine: Result %d", score)
				counter -= 1
				ChessLogger.Debug("Quiesce Goroutine: WG size: %d", counter)

				wg.Done()
				done <- ThreadedStatus{true, Mate + 1}
			}
			wg.Done()
		}(p)
	}

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

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("229: Result %d", maxValue)

	return maxValue
}

func Search(board ChessBoard, depth, alpha, beta, nodeType int, elderBrother bool) int {

	var res int

	if elderBrother {
		ChessLogger.Info("YBMW applied")
		res = LocalSearch(board, depth, alpha, beta, nodeType)
	} else {
		ChessLogger.Info("Offloading computation")
		res = RunRemote([]string{"Search", "\"" + ToEPD(board) + "\"", strconv.Itoa(depth), strconv.Itoa(alpha), strconv.Itoa(beta), strconv.Itoa(nodeType)})
	}
	return res
}
