package evaluation

import "github.com/mirceaIordache/goChess/common"

func EvaluateDraw(board common.ChessBoard) bool {
	common.ChessLogger.Info("Entering")
	common.ChessLogger.Debug("Board %s", common.ToEPD(board))
	/* Function that evaluates the draw score for given board */

	/* TODO Add some heuristics maybe, to speed it up? */

	whiteBoard := board.Board[common.White]
	blackBoard := board.Board[common.Black]

	if whiteBoard[common.Pawn] != 0 || blackBoard[common.Pawn] != 0 {
		common.ChessLogger.Info("Pawns on board. Exiting")
		return false
	}

	whiteMaterial := board.Material[common.White]
	blackMaterial := board.Material[common.Black]

	whiteKnights := common.NumBits(whiteBoard[common.Knight])
	blackKnights := common.NumBits(blackBoard[common.Knight])

	if (whiteMaterial < common.ValueRook || (whiteMaterial == 2*common.ValueKnight && whiteKnights == 2)) && (blackMaterial < common.ValueRook || (blackMaterial == 2*common.ValueKnight && blackKnights == 2)) {
		common.ChessLogger.Info("Not enough material. Exiting")
		return true
	}

	if whiteMaterial < common.ValueRook {
		if blackMaterial == 2*common.ValueBishop && (common.NumBits(blackBoard[common.Bishop]&common.WhiteSquares) == 2 || common.NumBits(blackBoard[common.Bishop]&common.BlackSquares) == 2) {
			common.ChessLogger.Info("One side has has Rook, other has 2 bishops on same colour. Exiting")
			return true
		}
	}

	if blackMaterial < common.ValueRook {
		if whiteMaterial == 2*common.ValueBishop && (common.NumBits(whiteBoard[common.Bishop]&common.WhiteSquares) == 2 || common.NumBits(whiteBoard[common.Bishop]&common.BlackSquares) == 2) {
			common.ChessLogger.Info("One side has has Rook, other has 2 bishops on same colour")
			return true
		}
	}

	common.ChessLogger.Info("Draw not detected. Exiting")
	return false
}
