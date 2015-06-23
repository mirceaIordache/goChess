package evaluation

import . "github.com/mirceaIordache/goChess/common"

func LoneKing(board ChessBoard, side, loser uint16) int {
	/* Lone King scenario */
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s, side %d, loser %d", ToEPD(board), side, loser)
	winner := loser ^ 1

	if board.Material[winner] == ValueBishop+ValueKnight && NumBits(board.Board[winner][Bishop]) == 1 && NumBits(board.Board[winner][Knight]) == 1 {
		ChessLogger.Info("KBNK Detected")
		res := ScoreKBNK(board, side, loser)
		ChessLogger.Debug("Result %d", res)
		return res
	}

	squareWinner := board.KingPos[winner]
	squareLoser := board.KingPos[loser]

	score := 150 - 6*TaxiDist[squareWinner][squareLoser] - EndingKing[squareLoser]
	if side == loser {
		score = -score
	}
	score += CalculateMaterial(board)

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result %d", score)
	return score
}

func ScoreKBNK(board ChessBoard, side, loser uint16) int {
	/* KBNK scoring. Based on GNU Chess implementation */
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s, side %d, loser %d", ToEPD(board), side, loser)
	winner := loser ^ 1
	squareB := board.KingPos[loser]
	if board.Board[winner][Bishop]&WhiteSquares != NullBitBoard {
		squareB = squareB&7*8 + 7 - squareB>>3
	}

	squareWinner := board.KingPos[winner]
	squareLoser := board.KingPos[loser]

	score := 300 - 6*TaxiDist[squareWinner][squareLoser]

	score -= KBNK[squareB]
	score -= EndingKing[squareLoser]
	score -= TaxiDist[LeadBit(board.Board[winner][Knight])][squareLoser]
	score -= TaxiDist[LeadBit(board.Board[winner][Bishop])][squareLoser]

	if board.Board[winner][King]&BitBoard(0x00003C3C3C3C0000) != NullBitBoard {
		score += 20
	}
	if side == loser {
		score = -score
	}

	score += CalculateMaterial(board)

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result %d", score)
	return score
}
