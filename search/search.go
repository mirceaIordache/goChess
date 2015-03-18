package search

import (
	"fmt"

	. "github.com/mirceaIordache/goChess/attack"
	. "github.com/mirceaIordache/goChess/common"
	. "github.com/mirceaIordache/goChess/evaluation"
	. "github.com/mirceaIordache/goChess/moveGenerator"
	. "github.com/mirceaIordache/goChess/quiesce"
	. "github.com/mirceaIordache/goChess/sort"
)

var overallBest = -Mate - 1

func SetBest(best int) {
	overallBest = best
}

func SearchRoot(board ChessBoard, iterationDepth, alpha, beta int) Move {
	fmt.Println("New Search: alpha: ", alpha, "  beta: ", beta)
	var score int
	side := board.Side
	best := -Mate
	nodeType := PV
	bestMove := Move{}
	list := GenerateMoves(board)
	FilterMoves(board, list, alpha, beta)
	iter := list

	for iter.Next != nil {
		p := PickBest(&iter)
		newBoard := ApplyMove(board, p.Move, side)
		if p == list.Value {
			score = -Search(newBoard, 2, iterationDepth-1, -beta, -alpha, nodeType)
			if beta == Mate && score <= alpha {
				alpha = -Mate
				score = -Search(board, 2, iterationDepth-1, -beta, -alpha, nodeType)
			}
		} else {
			nodeType = CUT
			if best > alpha {
				alpha = best
			}
			score = -Search(board, 2, iterationDepth-1, -alpha-1, -alpha, nodeType)
			if score > best {
				if alpha < score && score < beta {
					nodeType = PV
					score = -Search(board, 2, iterationDepth-1, -beta, -score, nodeType)
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
			return bestMove
		}
		iter = iter.Next
	}

	return bestMove
}

func Search(board ChessBoard, ply, depth, alpha, beta, nodeType int) int {
	if EvaluateDraw(board) == true || ply >= MaxPlyDepth {
		return SCORE_DRAW
	}

	list := GenerateMoves(board)
	FilterMoves(board, list, alpha, beta)

	var score int
	best := -Mate
	side := board.Side
	xside := 1 ^ side
	inCheck := SquareAttacked(board, board.KingPos[side], xside)
	if inCheck == true {
		list = GenerateCheckEscapes(board)
		if list.Next == nil {
			return Mate
		}
	}

	if overallBest >= Mate {
		return CalculateMaterial(board)
	}

	if depth <= 0 {
		return Quiesce(board, alpha, beta, 1).Score
	}

	if depth > 1 && nodeType != PV && inCheck == false && CalculateMaterial(board)+ValuePawn > beta && beta > -Mate && board.PawnMaterial[side] > ValueBishop {
		newBoard := ApplyNullMove(board, 0, side)
		nullScore := -Search(newBoard, ply+1, depth-3, -beta, -beta+1, nodeType)
		if nullScore >= beta {
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

	for true {
		p := PickBest(&list)
		newBoard := ApplyMove(board, p.Move, side)
		if firstMove {
			firstMove = false
			score = -Search(newBoard, ply+1, depth-1, -beta, -alpha, nodeType)
		} else {
			if list.Next == nil {
				break
			}
			if SquareAttacked(newBoard, newBoard.KingPos[side], xside) {
				continue
			}
			if nodeType == PV {
				nodeType = CUT
			}
			if best > alpha {
				alpha = best
			}
			score = -Search(newBoard, ply+1, depth-1, -alpha-1, -alpha, nodeType)
			if score > best {
				if saveNodeType == PV {
					nodeType = PV
				}
				if alpha < score && score < beta {
					score = -Search(newBoard, ply+1, depth-1, -beta, -score, nodeType)
				}
				if nodeType == PV && score <= alpha {
					score = -Search(newBoard, ply+1, depth-1, -alpha, -Mate, nodeType)
				}
			}
		}
		if score > best {
			best = score
			if best >= beta {
				break
			}
		}

		if Mate+1 == best {
			break
		}
	}

	return best
}
