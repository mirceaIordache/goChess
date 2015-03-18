package evaluation

import (
	"math"

	"github.com/mirceaIordache/goChess/attack"
	"github.com/mirceaIordache/goChess/common"
)

var maxPositionScore = [2]int{150, 150}
var lazyScore = [2]int{150, 150}

/* Evaluate the board */
/* Exclude possibility of draw, as it's checked by the EvaluateDraw function */
func Evaluate(alpha, beta int, board common.ChessBoard) int {
	/* Check for a mate in bounds */
	xside := 1 ^ board.Side
	if alpha > common.Mate-255 || beta < -common.Mate+255 {
		return int(board.Material[board.Side]) - int(board.Material[xside])
	}

	/* Check for KPK game */
	if common.CalculateMaterial(board) == common.ValuePawn {
		return KPK(board)
	}

	/* Check if there is a lone king and other side has no pawns */
	if board.Material[xside] == 0 && board.Board[board.Side][common.Pawn] == common.NullBitBoard {
		return LoneKing(board, board.Side, xside)
	}
	if board.Material[board.Side] == 0 && board.Board[xside][common.Pawn] == common.NullBitBoard {
		return LoneKing(board, xside, board.Side)
	}
	/* Attempt Lazy evaluation (calculate scores for bonuses/penalties) */

	var pieces [2]common.BitBoard
	var numPieces [2]uint16

	b := board.Board[common.White]
	pieces[common.White] = b[common.Knight] | b[common.Bishop] | b[common.Rook] | b[common.Queen]
	numPieces[common.White] = common.NumBits(pieces[common.White])

	b = board.Board[common.Black]
	pieces[common.Black] = b[common.Knight] | b[common.Bishop] | b[common.Rook] | b[common.Queen]
	numPieces[common.Black] = common.NumBits(pieces[common.Black])

	scoreTemp := common.CalculateMaterial(board)
	var score int
	if (scoreTemp+maxPositionScore[board.Side] < alpha || scoreTemp-maxPositionScore[xside] > beta) && common.GetGamePhase(board) <= 6 {
		/* Not end-game, so avoid lazy evaluation */
		score = scoreTemp
		goto next
	}

	score = 0
	score += ScoreDev(board, board.Side) - ScoreDev(board, xside)
	score += ScorePawn(board, board.Side) - ScorePawn(board, xside)
	score += ScoreKing(board, board.Side) - ScoreKing(board, xside)
	score += BishopTrapped(board, board.Side) - BishopTrapped(board, xside)
	score += DoubleQR7(board, board.Side) - DoubleQR7(board, xside)
	score += common.CalculateMaterial(board)

	/* Attempt a lazy evaluation cut */
	if (score+lazyScore[board.Side] < alpha || score-lazyScore[board.Side] > beta) == false {
		score = scoreTemp
	} else {
		attackBoard := attack.GenerateAttacks(board)
		scoreTemp = SCORE_HUNG * (HungPieces(board, board.Side, attackBoard) - HungPieces(board, xside, attackBoard))
		pinnedBoard := attack.FindPins(board)

		for piece := common.Knight; piece < common.King; piece++ {
			scoreTemp += ScorePiece(piece, board, board.Side, pinnedBoard) - ScorePiece(piece, board, xside, pinnedBoard)
		}

		lazyScore[board.Side] = int(math.Max(float64(scoreTemp), float64(lazyScore[board.Side])))
		maxPositionScore[board.Side] = int(math.Max(float64(scoreTemp), float64(maxPositionScore[board.Side])))
		score += scoreTemp + common.CalculateMaterial(board)
	}
next:
	/* Trade down bonus. When ahead, prefer pieces > pawns */
	if common.CalculateMaterial(board) >= 200 {
		/* Trade pieces score */
		score += int((common.NumBits(board.Friends[common.White]|board.Friends[common.Black]) - common.NumBits(pieces[common.White]|pieces[common.Black]))) * SCORE_TRADE_PIECE
		score -= int((common.NumBits(board.Board[common.White][common.Pawn]|board.Board[common.Black][common.Pawn]) - common.NumBits(board.Board[common.White][common.Pawn]|board.Board[common.Black][common.Pawn]))) * SCORE_TRADE_PAWN

	} else if common.CalculateMaterial(board) <= -200 {
		/* Trade pawns score */
		score -= int((common.NumBits(board.Friends[common.White]|board.Friends[common.Black]) - common.NumBits(pieces[common.White]|pieces[common.Black]))) * SCORE_TRADE_PIECE
		score += int((common.NumBits(board.Board[common.White][common.Pawn]|board.Board[common.Black][common.Pawn]) - common.NumBits(board.Board[common.White][common.Pawn]|board.Board[common.Black][common.Pawn]))) * SCORE_TRADE_PAWN
	}

	/* Heuristic, opposite colour bishops usually lead to draw.
	Adjust score accordingly */
	/* It ain't pretty, but it doesn't compile otherwise */
	if common.GetGamePhase(board) >= 6 && pieces[common.White] == board.Board[common.White][common.Bishop] && pieces[common.Black] == board.Board[common.Black][common.Bishop] && ((pieces[common.White]&common.WhiteSquares != 0 && pieces[common.Black]&common.BlackSquares != 0) || (pieces[common.White]&common.BlackSquares != 0 && pieces[common.Black]&common.WhiteSquares != 0)) {
		score = score / 2
	}

	/* Adjust score if a side has no mating material */
	if score > 0 && board.Board[board.Side][common.Pawn] == 0 && (board.Material[board.Side] < common.ValueRook ||
		pieces[board.Side] == board.Board[board.Side][common.Knight]) {
		score = 0
	}
	if score < 0 && board.Board[xside][common.Pawn] == 0 && (board.Material[xside] < common.ValueRook ||
		pieces[xside] == board.Board[xside][common.Knight]) {
		score = 0
	}

	return score
}
