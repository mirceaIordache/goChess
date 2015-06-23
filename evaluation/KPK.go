package evaluation

import . "github.com/mirceaIordache/goChess/common"

/* Evaluate a KPK endgame */
func KPK(board ChessBoard) int {
	var winner uint16

	if board.Board[White][Pawn] != NullBitBoard {
		winner = White
	} else {
		winner = Black
	}
	loser := 1 ^ winner

	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s, winner %d", ToEPD(board), winner)

	squarePawn := LeadBit(board.Board[winner][Pawn])
	squareWinner := board.KingPos[winner]
	squareLoser := board.KingPos[loser]

	delta := 0
	if winner == White {
		delta = int(squareWinner & 7)
	} else {
		delta = 7 - int(squareWinner&7)
	}
	score := ValuePawn + (ValueQueen * passedPawnScore[winner][squarePawn&7] / SCORE_PFACTOR) + 4*delta
	/* Check if pawn is outside the square of the king */

	if ^SquarePawnMask[winner][squarePawn]&board.Board[loser][King] != NullBitBoard {
		if MoveArray[King][squareLoser]&SquarePawnMask[winner][squarePawn] == NullBitBoard {
			if winner == board.Side {
				ChessLogger.Info("Loser King too far away from pawn. Exiting")
				ChessLogger.Debug("Result %d", score)
				return score
			} else {
				ChessLogger.Info("Loser King too far away from pawn. Exiting")
				ChessLogger.Debug("Result %d", -score)
				return -score
			}
		}
		if winner == board.Side {
			ChessLogger.Info("Exiting")
			ChessLogger.Debug("Result %d", score)
			return score
		}
	}

	/* Friendly king is on same or adjacent file to pawn;
	Pawn on any other file than rook (A, H) */
	if squarePawn>>3 != 0 && squarePawn>>3 != 7 && (IsolaniMask[squarePawn>>3]|FileBit[squarePawn]>>3)&board.Board[winner][King] != NullBitBoard {

		/* Case check:
		1. Friendly king is 2 ranks more advanced than pawn
		2. Friendly king is 1 rank more advanced than pawn:
			a. Friendly king is on 6th rank
			b. Enemy king doesn't have direct opposition by being 2 ranks in front of friendly king, and on same file
		3. Friendly king is same rank as pawn:
			a. Enemy king is not 2-4 ranks more advanced than pawn
			b. Pawn is on 6th rank and enemy king does not have direct opposition
		4. Pawn on 7th rank, friendly king on 6th rank
			a. Enemy king not on queening square
			b. Enemy king is on queening square, but both kings on same file
		*/

		if winner == White {
			if squareWinner&7 == (squarePawn&7 + 2) {
				if winner == board.Side {
					ChessLogger.Info("Friendly king is 2 ranks ahead of pawn. Exiting")
					ChessLogger.Debug("Result %d", score)
					return score
				} else {
					ChessLogger.Info("Friendly king is 2 ranks ahead of pawn. Exiting")
					ChessLogger.Debug("Result %d", -score)
					return -score
				}
			}
			if squareWinner&7 == (squarePawn&7 + 1) {
				if squareWinner&7 == 5 {
					if winner == board.Side {
						ChessLogger.Info("Friendly king is 1 rank ahead of pawn on 6th rank. Exiting")
						ChessLogger.Debug("Result %d", score)
						return score
					} else {
						ChessLogger.Info("Friendly king is 1 rank ahead of pawn on 6th rank. Exiting")
						ChessLogger.Debug("Result %d", -score)
						return -score
					}
				}
				if squareWinner < A6 {
					if squareWinner+16 == squareLoser && winner == board.Side {
						ChessLogger.Info("Friendly king is 1 rank ahead of pawn. Enemy king 1 rank behind. Exiting")
						ChessLogger.Debug("Result %d", 0)
						return 0
					} else {
						if winner == board.Side {
							ChessLogger.Info("Friendly king is 1 rank ahead of pawn. Enemy king 1 rank behind. Exiting")
							ChessLogger.Debug("Result %d", score)
							return score
						} else {
							ChessLogger.Info("Friendly king is 1 rank ahead of pawn. Enemy king 1 rank behind. Exiting")
							ChessLogger.Debug("Result %d", -score)
							return -score
						}
					}
				}
			}

			if squareWinner&7 == squarePawn&7 {
				if (squareLoser&7-squarePawn&7 < 2 || squareLoser&7-squarePawn&7 > 4) && winner == board.Side {
					ChessLogger.Info("Friendly king on the same rank as pawn. Enemy king out of range")
					ChessLogger.Debug("Result %d", score)
					return score
				}
				if (squareLoser&7-squarePawn&7 < 1 || squareLoser&7-squarePawn&7 > 5) && loser == board.Side {
					ChessLogger.Info("Friendly king on the same rank as pawn. Enemy king out of range")
					ChessLogger.Debug("Result %d", -score)

					return -score
				}
				if squarePawn&7 == 5 && squareWinner+16 != squareLoser {
					if winner == board.Side {
						ChessLogger.Info("Pawn on 6th rank. Friendly king on same rank, 2 files ahead of enemy king")
						ChessLogger.Debug("Result %d", score)
						return score
					} else {

						ChessLogger.Info("Pawn on 6th rank. Friendly king on same rank, 2 files ahead of enemy king")
						ChessLogger.Debug("Result %d", 0)
						return 0
					}
				}

				if squarePawn&7 == 6 && squareWinner&7 == 5 {
					if squareLoser != squarePawn+8 {
						if winner == board.Side {
							ChessLogger.Info("Pawn on 7th rank, friendly king on 6th rank. Enemy king is 1 file behind pawn")
							ChessLogger.Debug("Result %d", score)
							return score
						} else {

							ChessLogger.Info("Pawn on 7th rank, friendly king on 6th rank. Enemy king is 1 file behind pawn")
							ChessLogger.Debug("Result %d", 0)
							return 0
						}
					}
					if squareLoser == squarePawn+8 && squareLoser == squareWinner+16 {
						if winner == board.Side {
							ChessLogger.Info("Enemy king 1 file ahead of pawn. Friendly king on 2 files behind of enemy king.")
							ChessLogger.Debug("Result %d", score)

							return score
						} else {
							ChessLogger.Info("Enemy king 1 file ahead of pawn. Friendly king on 2 files behind of enemy king.")
							ChessLogger.Debug("Result %d", 0)

							return 0
						}
					}
				}
			}
		} else {
			if squareWinner&7 == squarePawn&7-2 {
				if winner == board.Side {
					ChessLogger.Info("Friendly king is 2 ranks ahead of pawn")
					ChessLogger.Debug("Result %d", score)

					return score
				} else {

					ChessLogger.Info("Friendly king is 2 ranks ahead of pawn")
					ChessLogger.Debug("Result %d", -score)
					return -score
				}
			}
			if squareWinner&7 == squarePawn&7-1 {
				if squareWinner&7 == 2 {
					if winner == board.Side {

						ChessLogger.Info("Friendly king is 1 rank ahead of pawn")
						ChessLogger.Debug("Result %d", score)
						return score
					} else {

						ChessLogger.Info("Friendly king is 1 rank ahead of pawn")
						ChessLogger.Debug("Result %d", -score)
						return -score
					}
				}
				if squareWinner > H3 {
					if squareWinner-16 == squareLoser && board.Side == winner {

						ChessLogger.Info("Friendly king 2 files behind enemy king. Pawn approaching promotion")
						ChessLogger.Debug("Result %d", 0)
						return 0
					} else {
						if winner == board.Side {

							ChessLogger.Info("Pawn approaching promotion")
							ChessLogger.Debug("Result %d", score)
							return score
						} else {

							ChessLogger.Info("Pawn approaching promotion")
							ChessLogger.Debug("Result %d", -score)
							return -score
						}
					}
				}
			}
			if squareWinner&7 == squarePawn&7 {
				if squarePawn&7-squareLoser&7 < 2 || squarePawn&7-squareLoser&7 > 4 && winner == board.Side {

					ChessLogger.Info("Friendly king on same rank as pawn. Enemy king is outside pawn square")
					ChessLogger.Debug("Result %d", score)
					return score
				}
				if squarePawn&7-squareLoser&7 < 1 || squarePawn&7-squareLoser&7 > 5 && loser == board.Side {
					ChessLogger.Info("Friendly king on same rank as pawn. Enemy king is outside pawn square")
					ChessLogger.Debug("Result %d", -score)
					return -score
				}
				if squarePawn&7 == 5 && squareWinner+16 != squareLoser {
					if winner == board.Side {
						ChessLogger.Info("Friendly king and pawn on 6th rank. Enemy king is 2 files ahead of king")
						ChessLogger.Debug("Result %d", score)
						return score
					} else {
						ChessLogger.Info("Friendly king and pawn on 6th rank. Enemy king is 2 files ahead of king")
						ChessLogger.Debug("Result %d", -score)
						return 0
					}
				}
			}

			if squarePawn&7 == 1 && squareWinner&7 == 2 {
				if squareLoser != squarePawn-8 {
					if winner == board.Side {
						ChessLogger.Info("Pawn on 2nd rank, friendly king on 3rd rank. Enemy king is not in front of pawn")
						ChessLogger.Debug("Result %d", score)
						return score
					} else {
						ChessLogger.Info("Pawn on 2nd rank, friendly king on 3rd rank. Enemy king is not in front of pawn")
						ChessLogger.Debug("Result %d", -score)
						return 0
					}
				}
				if squareLoser == squarePawn-8 && squareLoser == squareWinner-16 {
					if winner == board.Side {
						ChessLogger.Info("Pawn on 2nd rank. Enemy king is in front of pawn and friendly king is behind pawn")
						ChessLogger.Debug("Result %d", score)
						return score
					} else {
						ChessLogger.Info("Pawn on 2nd rank. Enemy king is in front of pawn and friendly king is behind pawn")
						ChessLogger.Debug("Result %d", 0)
						return 0
					}
				}
			}
		}
	}
	ChessLogger.Info("Inconclusive situation")
	ChessLogger.Debug("Result %d", 0)
	return 0
}
