package evaluation

import "github.com/mirceaIordache/goChess/common"

func EvaluateDraw(board common.ChessBoard) bool {
	/* Function that evaluates the draw score for given board */

	/* TODO Add some heuristics maybe, to speed it up? */

	whiteBoard := board.Board[common.White]
	blackBoard := board.Board[common.Black]

	if whiteBoard[common.Pawn] != 0 || blackBoard[common.Pawn] != 0 {
		return false
	}

	whiteMaterial := board.Material[common.White]
	blackMaterial := board.Material[common.Black]

	whiteKnights := common.NumBits(whiteBoard[common.Knight])
	blackKnights := common.NumBits(blackBoard[common.Knight])

	if (whiteMaterial < common.ValueRook || (whiteMaterial == 2*common.ValueKnight && whiteKnights == 2)) && (blackMaterial < common.ValueRook || (blackMaterial == 2*common.ValueKnight && blackKnights == 2)) {
		return true
	}

	if whiteMaterial < common.ValueRook {
		if blackMaterial == 2*common.ValueBishop && (common.NumBits(blackBoard[common.Bishop]&common.WhiteSquares) == 2 || common.NumBits(blackBoard[common.Bishop]&common.BlackSquares) == 2) {
			return true
		}
	}

	if blackMaterial < common.ValueRook {
		if whiteMaterial == 2*common.ValueBishop && (common.NumBits(whiteBoard[common.Bishop]&common.WhiteSquares) == 2 || common.NumBits(whiteBoard[common.Bishop]&common.BlackSquares) == 2) {
			return true
		}
	}

	return false
}
