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
				return score
			} else {
				return -score
			}
		}
		if winner == board.Side {
			return score
		}
	}

	/* Friendly king is on same or adjacent file to pawn;
	Pawn on any other file than rook (A, H) */
	if squarePawn>>3 != 0 && squarePawn>>3 != 7 && (IsolaniMask[squarePawn>>3]|FileBit[squarePawn>>3])&board.Board[winner][King] != NullBitBoard {

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
					return score
				} else {
					return -score
				}
			}
			if squareWinner&7 == (squarePawn&7 + 1) {
				if squareWinner&7 == 5 {
					if winner == board.Side {
						return score
					} else {
						return -score
					}
				}
				if squareWinner < A6 {
					if squareWinner+16 == squareLoser && winner == board.Side {
						return 0
					} else {
						if winner == board.Side {
							return score
						} else {
							return -score
						}
					}
				}
			}

			if squareWinner&7 == squarePawn&7 {
				if (squareLoser&7-squarePawn&7 < 2 || squareLoser&7-squarePawn&7 > 4) && winner == board.Side {
					return score
				}
				if (squareLoser&7-squarePawn&7 < 1 || squareLoser&7-squarePawn&7 > 5) && loser == board.Side {
					return -score
				}
				if squarePawn&7 == 5 && squareWinner+16 != squareLoser {
					if winner == board.Side {
						return score
					} else {
						return 0
					}
				}

				if squarePawn&7 == 6 && squareWinner&7 == 5 {
					if squareLoser != squarePawn+8 {
						if winner == board.Side {
							return score
						} else {
							return 0
						}
					}
					if squareLoser == squarePawn+8 && squareLoser == squareWinner+16 {
						if winner == board.Side {
							return score
						} else {
							return 0
						}
					}
				}
			}
		} else {
			if squareWinner&7 == squarePawn&7-2 {
				if winner == board.Side {
					return score
				} else {
					return -score
				}
			}
			if squareWinner&7 == squarePawn&7-1 {
				if squareWinner&7 == 2 {
					if winner == board.Side {
						return score
					} else {
						return -score
					}
				}
				if squareWinner > H3 {
					if squareWinner-16 == squareLoser && board.Side == winner {
						return 0
					} else {
						if winner == board.Side {
							return score
						} else {
							return -score
						}
					}
				}
			}
			if squareWinner&7 == squarePawn&7 {
				if squarePawn&7-squareLoser&7 < 2 || squarePawn&7-squareLoser&7 > 4 && winner == board.Side {
					return score
				}
				if squarePawn&7-squareLoser&7 < 1 || squarePawn&7-squareLoser&7 > 5 && loser == board.Side {
					return -score
				}
				if squarePawn&7 == 5 && squareWinner+16 != squareLoser {
					if winner == board.Side {
						return score
					} else {
						return 0
					}
				}
			}

			if squarePawn&7 == 1 && squareWinner&7 == 2 {
				if squareLoser != squarePawn-8 {
					if winner == board.Side {
						return score
					} else {
						return 0
					}
				}
				if squareLoser == squarePawn-8 && squareLoser == squareWinner-16 {
					if winner == board.Side {
						return score
					} else {
						return 0
					}
				}
			}
		}
	}
	return 0
}
