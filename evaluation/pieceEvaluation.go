package evaluation

import (
	"math"

	. "github.com/mirceaIordache/goChess/attack"
	. "github.com/mirceaIordache/goChess/common"
)

func ScoreDev(board ChessBoard, side uint16) int {
	/* Development potential score

	Penalise: 	Uncastled and cannot be castled
				Undeveloped Knights and Bishops
				Early Queen move
	*/

	score := 0
	moveBoard := GenerateMoveBoard(board)
	developed := (board.Board[side][Knight] & nn[side]) | (board.Board[side][Bishop] & bb[side])
	score = int(NumBits(developed)) * -8

	if board.Castled[side] == true && board.GameCount >= 38 {
		return score
	}

	score += SCORE_DEV_NOTCASTLED

	if moveBoard[board.KingPos[side]] > 0 {
		score += SCORE_DEV_KINGMOVED
	}

	ourRooks := board.Board[side][Rook]

	for ourRooks != NullBitBoard {
		sq := LeadBit(ourRooks)
		ClearBit(&ourRooks, int(sq))
		if moveBoard[sq] > 0 {
			score += SCORE_DEV_ROOKMOVED
		}
	}

	if board.Board[side][Queen] != NullBitBoard {
		sq := LeadBit(board.Board[side][Queen])
		if moveBoard[sq] > 0 {
			score += SCORE_QUEEN_EARLYMOVE
		}
	}

	developed = board.Board[side][Pawn] | BitBoard(0xc3c3c3c3c3c3c3c3)
	for developed != NullBitBoard {
		sq := LeadBit(developed)
		ClearBit(&developed, int(sq))

		if moveBoard[sq] > 0 {
			score += SCORE_PAWN_EARLYWINGMOVE
		}
	}

	return score
}

func ScorePawn(board ChessBoard, side uint16) int {
	/* NOTE: GNU Chess uses a hash table here for optimisation.
	Ignoring it for now. MMV
	*/

	/* Pawn Evaluation score

	Factors:	Passed pawns
				Backward pawns
				Pawn base attacked
				Doubled pawns
				Isolated pawns
				Connected passed pawns on 6th (or 7th) rank
				Unmoved, blocked D, E pawn
				Uncatchable passed pawn
				Pawn storms
	*/

	if board.Board[side][Pawn] == NullBitBoard {
		return 0
	}

	var nFile [8]int
	score := 0
	xside := 1 ^ side
	enemyKing := board.KingPos[xside]
	enemyPawns := board.Board[xside][Pawn]
	ourPawns := board.Board[side][Pawn]
	throwAway := board.Board[side][Pawn]
	passed := NullBitBoard
	weaked := NullBitBoard

	for throwAway != 0 {
		sq := LeadBit(throwAway)
		ClearBit(&throwAway, int(sq))
		score += pawnScores[side][sq]

		/* Use raycasting to find if passed pawn */
		if (enemyPawns & PassedPawnMask[side][sq]) == NullBitBoard {
			if (side == White && (FromToRay[sq][sq|56]&ourPawns) == 0) || (side == Black && (FromToRay[sq][sq&7]&ourPawns) == 0) {
				passed |= BitPosArray[sq]
				score += (passedPawnScore[side][sq>>3] * int(GetGamePhase(board))) / 12
			}
		}

		var i uint16
		backward := false
		if side == White {
			i = sq + 8
		} else {
			i = sq - 8
		}

		if (PassedPawnMask[xside][i] & ^FileBit[sq&7] & ourPawns) == 0 && GenerateCBoard(board)[i] != Pawn {
			n1 := NumBits(ourPawns & MoveArray[PawnType[xside]][i])
			n2 := NumBits(enemyPawns & MoveArray[PawnType[side]][i])
			if n1 < n2 {
				backward = true
			}
		}

		if backward == false && (BitPosArray[sq]&brank7[xside]) != 0 {
			i1 := 1
			if side == White {
				i += 8
			} else {
				i -= 8
			}
			if (PassedPawnMask[xside][i] & ^FileBit[i1&7] & ourPawns) == 0 {
				n1 := NumBits(ourPawns & MoveArray[PawnType[xside]][i])
				n2 := NumBits(enemyPawns & MoveArray[PawnType[xside]][i])
				if n1 < n2 {
					backward = true
				}
			}
		}

		if backward == true {
			weaked |= BitPosArray[sq]
			score += -(8 + int(GetGamePhase(board)))
		}

		if MoveArray[PawnType[side]][sq]&enemyPawns != 0 && MoveArray[PawnType[side]][sq]&ourPawns != 0 {
			score += SCORE_PAWN_BASEATTACK
		}

		nFile[sq&7]++
	}

	for i := 0; i <= 7; i++ {
		if nFile[i] > 1 {
			score += -(8 + int(GetGamePhase(board)))
		}

		if nFile[i] != 0 && ourPawns&IsolaniMask[i] == 0 {
			if FileBit[i]&enemyPawns == 0 {
				score += IsolaniWeaker[i] * nFile[i]
			} else {
				score += IsolaniNormal[i] * nFile[i]
			}
			weaked |= ourPawns & FileBit[i]
		}
	}

	if board.Side == board.OurColor {
		if NumBits(enemyPawns) == 8 {
			score += SCORE_PAWN_EIGHT
		}

		if NumBits(Stonewall[xside]&enemyPawns) == 3 {
			score += SCORE_PAWN_STONEWALL
		}

		mobility := 0
		if side == White {
			mobility = int(NumBits((ourPawns >> 8) & enemyPawns & Boxes[1]))
		} else {
			mobility = int(NumBits((ourPawns << 8) & enemyPawns & Boxes[1]))
		}

		if mobility > 1 {
			score += mobility * SCORE_PAWN_LOCKED
		}
	}

	if side == White && board.Board[side][Queen] != 0 && (BitPosArray[C6]|BitPosArray[F6])&ourPawns != 0 {
		if ourPawns&BitPosArray[F6] != 0 && enemyKing > H6 && Distance[enemyKing][G7] == 1 {
			score += SCORE_PAWN_NEARKING
		}
		if ourPawns&BitPosArray[C6] != 0 && enemyKing > H6 && Distance[enemyKing][B7] == 1 {
			score += SCORE_PAWN_NEARKING
		}
	} else if side == Black && board.Board[side][Queen] != 0 && (BitPosArray[C3]|BitPosArray[F3])&ourPawns != 0 {
		if ourPawns&BitPosArray[F3] != 0 && enemyKing < A3 && Distance[enemyKing][G2] == 1 {
			score += SCORE_PAWN_NEARKING
		}
		if ourPawns&BitPosArray[C3] != 0 && enemyKing < A3 && Distance[enemyKing][B2] == 1 {
			score += SCORE_PAWN_NEARKING
		}
	}

	throwAway = passed & brank67[side]
	if throwAway != 0 && (board.PawnMaterial[xside] == ValueRook || (board.PawnMaterial[xside] == ValueKnight && CalculatePieces(board, xside) == board.Board[xside][Knight])) {
		row := board.KingPos[xside] & 7
		rank := board.KingPos[xside] >> 3
		for i := 0; i <= 6; i++ {
			if throwAway&FileBit[i] != 0 && throwAway&FileBit[i+1] != 0 && (int(row) < i-1 || int(row) > i+1 || (side == White && rank < 4) || (side == Black && rank > 3)) {
				score += SCORE_PAWN_CONNECTED
			}
		}
	}

	blocker := board.Friends[side] | board.Friends[xside]
	if side == White && (((ourPawns*d2e2[White])>>8)&blocker != 0) {
		score += SCORE_PAWN_BLOCKED
	}
	if side == Black && (((ourPawns*d2e2[Black])<<8)&blocker != 0) {
		score += SCORE_PAWN_BLOCKED
	}

	if passed != 0 && board.PawnMaterial[xside] == 0 {
		//		enemy := board.Board[xside]
		throwAway = passed
		for throwAway != 0 {
			sq := LeadBit(throwAway)
			ClearBit(&throwAway, int(sq))
			if board.Side == side {
				if SquarePawnMask[side][sq]&board.Board[xside][King] == 0 {
					score += ValueQueen * passedPawnScore[side][sq>>3] / SCORE_PFACTOR
				}
			} else if MoveArray[King][enemyKing]&SquarePawnMask[side][sq] == 0 {
				score += ValueQueen * passedPawnScore[side][sq>>3] / SCORE_PFACTOR
			}
		}
	}

	if math.Abs(float64(board.KingPos[side]&7-board.KingPos[xside]&7)) >= 4 && GetGamePhase(board) < 6 {
		row := enemyKing & 7
		throwAway = (IsolaniMask[row] | FileBit[row]) & ourPawns
		for throwAway != 0 {
			sq := LeadBit(throwAway)
			ClearBit(&throwAway, int(sq))
			score += 10 * (5 - Distance[sq][enemyKing])
		}
	}

	return score
}

func ScoreKnight(board ChessBoard, side uint16, pinnedBoard BitBoard) int {
	/* Knight Evaluation score

	   Factors:	Central knight
				Mobility/Control/Attack
				Outpost knight
				Weak pawn attack (Maybe later?)
	*/

	if board.Board[side][Knight] == NullBitBoard {
		return 0
	}

	xside := side ^ 1
	score := 0
	scoreTemp := 0
	ourKnights := board.Board[side][Knight]
	enemyPawns := board.Board[xside][Pawn]

	if ourKnights&pinnedBoard != 0 {
		score += SCORE_KNIGHT_PINNED * int(NumBits(ourKnights&pinnedBoard))
	}

	for ourKnights != 0 {
		sq := LeadBit(ourKnights)
		ClearBit(&ourKnights, int(sq))

		scoreTemp = ScoreControl(board, sq, side)
		if (BitPosArray[sq] & Rings[3]) != NullBitBoard {
			scoreTemp += SCORE_KNIGHT_ONRIM
		}

		if Outpost[side][sq] != 0 && enemyPawns&IsolaniMask[sq&7]&PassedPawnMask[side][sq] == 0 {
			scoreTemp += SCORE_KNIGHT_OUTPOST
			if MoveArray[PawnType[xside]][sq]&board.Board[side][Pawn] != 0 {
				scoreTemp += SCORE_KNIGHT_OUTPOST
			}
		}

		score += scoreTemp

	}

	return score
}

func ScoreBishop(board ChessBoard, side uint16, pinnedBoard BitBoard) int {
	/* Bishop Evaluation Score

	   Factors:	Double Bishops
				Mobility/control/attack
				Outpost Bishop
				Fianchetto Bishop
				Bishop Pair
	*/

	if board.Board[side][Bishop] == NullBitBoard {
		return 0
	}

	xside := 1 ^ side
	score := 0
	scoreTemp := 0
	n := 0
	ourBishops := board.Board[side][Bishop]
	enemyPawns := board.Board[xside][Pawn]

	if ourBishops&pinnedBoard != 0 {
		score += SCORE_BISHOP_PINNED * int(NumBits(ourBishops&pinnedBoard))
	}

	for ourBishops != 0 {
		sq := LeadBit(ourBishops)
		ClearBit(&ourBishops, int(sq))
		n++

		scoreTemp = ScoreControl(board, sq, side)

		if Outpost[side][sq] != 0 && (enemyPawns&IsolaniMask[sq&7]&PassedPawnMask[side][sq]) == 0 {
			scoreTemp += SCORE_BISHOP_OUTPOST
			if MoveArray[PawnType[xside]][sq]&board.Board[side][Pawn] != 0 {
				scoreTemp += SCORE_BISHOP_OUTPOST
			}
		}
		if side == White {
			if board.KingPos[side] >= F1 && board.KingPos[side] <= H1 && sq == G2 {
				scoreTemp += SCORE_BISHOP_FIANCHETTO
			}
			if board.KingPos[side] >= A1 && board.KingPos[side] <= C1 && sq == B2 {
				scoreTemp += SCORE_BISHOP_FIANCHETTO
			}
		} else if side == Black {
			if board.KingPos[side] >= F8 && board.KingPos[side] <= H8 && sq == G7 {
				scoreTemp += SCORE_BISHOP_FIANCHETTO
			}
			if board.KingPos[side] >= A8 && board.KingPos[side] <= C8 && sq == B7 {
				scoreTemp += SCORE_BISHOP_FIANCHETTO
			}
		}

		score += scoreTemp
	}

	if n > 1 {
		score += SCORE_BISHOP_DOUBLE
	}
	return score
}

func ScoreRook(board ChessBoard, side uint16, pinnedBoard BitBoard) int {
	/* Rook Evaluation Score

	   Factors:	Rook on 7th Rank and Enemy King on 8th Rank or Pawns on 7th Rank
				Rook on open/half-open file
				Rook in front/behind passed pawns
	*/

	if board.Board[side][Rook] == NullBitBoard {
		return 0
	}

	score := 0
	scoreTemp := 0
	xside := 1 ^ side
	ourRooks := board.Board[side][Rook]
	enemyKing := board.KingPos[xside]

	if ourRooks&pinnedBoard != 0 {
		score += SCORE_ROOK_PINNED * int(NumBits(ourRooks&pinnedBoard))
	}

	for ourRooks != 0 {
		sq := LeadBit(ourRooks)
		ClearBit(&ourRooks, int(sq))

		scoreTemp = ScoreControl(board, sq, side)

		file := sq & 7
		if GetGamePhase(board) < 7 {
			if board.Board[side][Pawn]&FileBit[file] == 0 {
				if file == 5 && enemyKing&7 >= EFile {
					scoreTemp += SCORE_ROOK_LIBERATED
					scoreTemp += SCORE_ROOK_HALFFILE
					if board.Board[xside][Pawn]&FileBit[file] == 0 {
						scoreTemp += SCORE_ROOK_OPENFILE
					}
				}
			}
		}

		if sq>>3 == Rank7[side] && (enemyKing>>3 == Rank8[side] || board.Board[xside][Pawn]&RankBit[sq>>3] != 0) {
			scoreTemp += SCORE_ROOK_7RANK
		}

		score += scoreTemp
	}

	return score
}

func ScoreQueen(board ChessBoard, side uint16, pinnedBoard BitBoard) int {
	/* Queen Evaluation Score */

	if board.Board[side][Queen] == NullBitBoard {
		if side == board.OurColor {
			return SCORE_QUEEN_ABSENT
		}
		return 0
	}

	score := 0
	scoreTemp := 0
	xside := side ^ 1
	ourQueens := board.Board[side][Queen]
	enemyKing := board.KingPos[xside]

	if ourQueens&pinnedBoard != 0 {
		score += SCORE_QUEEN_PINNED * int(NumBits(ourQueens&pinnedBoard))
	}

	for ourQueens != 0 {
		sq := LeadBit(ourQueens)
		ClearBit(&ourQueens, int(sq))

		scoreTemp = ScoreControl(board, sq, side)

		if Distance[sq][enemyKing] <= 2 {
			scoreTemp += SCORE_QUEEN_NEARKING
		}

		score += scoreTemp
	}

	return score
}

func ScoreKing(board ChessBoard, side uint16) int {
	/* King Evaluation score

	Factors:	King in Corner
				Pawns around king
				King on open file
				Uncastled King
				Open Rook file via Ng5 or Bxh7 sac
				No Major or Minor piece in king's quadrant

	*/
	score := 0
	xside := 1 ^ side
	sq := board.KingPos[side]
	file := sq & 7
	rank := sq >> 3
	fianchettoSq := -1
	phase := GetGamePhase(board)
	moveBoard := GenerateMoveBoard(board)
	cboard := GenerateCBoard(board)

	if phase < 6 {
		score += (int(6-phase)*KingSq[sq] + int(phase)*EndingKing[sq]) / 6
		n := 0
		if side == White {
			n = int(NumBits(MoveArray[King][sq] & board.Board[side][Pawn] & RankBit[rank+1]))
		} else {
			n = int(NumBits(MoveArray[King][sq] & board.Board[side][Pawn] & RankBit[rank-1]))
		}
		score += pawnCover[n]
		if board.Castled[side] == false {
			n = -1
			if side == White {
				if sq == 4 && moveBoard[sq] == 0 {
					if (board.Board[side][Rook]&BitPosArray[H1]) != NullBitBoard && moveBoard[H1] == 0 {
						n = int(NumBits(MoveArray[King][G1] & board.Board[side][Pawn] & RankBit[rank+1]))
					}
					if (board.Board[side][Rook]&BitPosArray[A1]) != NullBitBoard && moveBoard[A1] == 0 {
						n = int(NumBits(MoveArray[King][C1] & board.Board[side][Pawn] & RankBit[rank+1]))
					}
				}
			} else {
				if sq == 60 && moveBoard[sq] == 0 {
					if (board.Board[side][Rook]&BitPosArray[H8]) != NullBitBoard && moveBoard[H8] == 0 {
						n = int(NumBits(MoveArray[King][G8] & board.Board[side][Pawn] & RankBit[rank-1]))
					}
					if (board.Board[side][Rook]&BitPosArray[A8]) != NullBitBoard && moveBoard[A8] == 0 {
						n = int(NumBits(MoveArray[King][C8] & board.Board[side][Pawn] & RankBit[rank-1]))
					}
				}
			}

			if n != -1 {
				score += pawnCover[n]
			}
		}
		if side == board.OurColor && file >= FFile && (FileBit[GFile]&board.Board[side][Pawn]) == 0 {
			if side == White && cboard[F2] == Pawn {
				score += SCORE_KING_GOPEN
			} else if side == Black && cboard[F7] == Pawn {
				score += SCORE_KING_GOPEN
			}
		}

		if FileBit[file]&board.Board[side][Pawn] == 0 {
			score += SCORE_KING_OPENFILE
		}

		if FileBit[file]&board.Board[xside][Pawn] == 0 {
			score += SCORE_KING_ENEMYOPENFILE
		}

		switch file {
		case AFile, EFile, FFile, GFile:
			{
				if (FileBit[file+1] & board.Board[side][Pawn]) == 0 {
					score += SCORE_KING_OPENFILE
				}
				if (FileBit[file+1] & board.Board[xside][Pawn]) == 0 {
					score += SCORE_KING_ENEMYOPENFILE
				}
			}
		case BFile, CFile, DFile, HFile:
			{
				if (FileBit[file-1] & board.Board[side][Pawn]) == 0 {
					score += SCORE_KING_OPENFILE
				}
				if (FileBit[file-1] & board.Board[xside][Pawn]) == 0 {
					score += SCORE_KING_ENEMYOPENFILE
				}
			}
		}

		if board.Castled[side] == true {
			if side == White {
				if file > EFile {
					if (BitPosArray[F2]&board.Board[side][Pawn]) == 0 || (BitPosArray[G2]&board.Board[side][Pawn]) == 0 || (BitPosArray[G2]&board.Board[side][Pawn]) == 0 {
						score += SCORE_KING_RUPTURE
					}
				} else if file < EFile {
					if (BitPosArray[A2]&board.Board[side][Pawn]) == 0 || (BitPosArray[B2]&board.Board[side][Pawn]) == 0 || (BitPosArray[C2]&board.Board[side][Pawn]) == 0 {
						score += SCORE_KING_RUPTURE
					}
				}
			} else {
				if file > EFile {
					if (BitPosArray[F7]&board.Board[side][Pawn]) == 0 || (BitPosArray[G7]&board.Board[side][Pawn]) == 0 || (BitPosArray[H7]&board.Board[side][Pawn]) == 0 {
						score += SCORE_KING_RUPTURE
					}
				} else if file < EFile {
					if (BitPosArray[A7]&board.Board[side][Pawn]) == 0 || (BitPosArray[B7]&board.Board[side][Pawn]) == 0 || (BitPosArray[C7]&board.Board[side][Pawn]) == 0 {
						score += SCORE_KING_RUPTURE
					}
				}
			}
		}

		/* 1:18 AM. First Energy Drink of the night */
		if side == board.OurColor {
			if file >= EFile && board.Board[xside][Queen] != NullBitBoard && board.Board[xside][Rook] != NullBitBoard && (board.Board[side][Pawn]|board.Board[xside][Pawn])&FileBit[7] == NullBitBoard {
				score += SCORE_KING_HOPEN
			}

			if side == White {
				if file > EFile {
					if board.Board[side][Rook]&MaskKR_TrappedWhite[HFile-file] != 0 {
						score += SCORE_ROOK_TRAPPED
					}
				} else if file < DFile {
					if board.Board[side][Rook]&MaskQR_TrappedWhite[file] != 0 {
						score += SCORE_ROOK_TRAPPED
					}
				}
			} else {
				if file > EFile {
					if board.Board[side][Rook]&MaskKR_TrappedBlack[HFile-file] != 0 {
						score += SCORE_ROOK_TRAPPED
					}
				} else if file < DFile {
					if board.Board[side][Rook]&MaskQR_TrappedBlack[file] != 0 {
						score += SCORE_ROOK_TRAPPED
					}
				}
			}
		}
		if file > EFile && board.KingPos[xside]&7 < DFile {
			if side == White {
				fianchettoSq = G3
			} else {
				fianchettoSq = G6
			}

			if (BitPosArray[fianchettoSq] & board.Board[side][Pawn]) != NullBitBoard {
				if (BitPosArray[F4]|BitPosArray[H4]|BitPosArray[F5]|BitPosArray[H5])&board.Board[xside][Pawn] != NullBitBoard {
					score += SCORE_PAWN_FIANCHETTO
				}
			}
		}

		if file < EFile && board.KingPos[xside]&7 > EFile {
			if side == White {
				fianchettoSq = B3
			} else {
				fianchettoSq = B6
			}

			if (BitPosArray[fianchettoSq] & board.Board[side][Pawn]) != NullBitBoard {
				if (BitPosArray[A4]|BitPosArray[C4]|BitPosArray[A5]|BitPosArray[C5])&board.Board[xside][Pawn] != NullBitBoard {
					score += SCORE_PAWN_FIANCHETTO
				}
			}
		}

		var i int
		if file <= DFile {
			i = 1
		} else {
			i = 0
		}

		x := BoardHalf[side] & BoardSide[i]
		n1 := NumBits(x & board.Friends[xside])
		if n1 > 0 {
			n2 := NumBits(x & (board.Friends[side] & ^board.Board[side][Pawn] & ^board.Board[side][King]))
			if n1 > n2 {
				score += int((n1 - n2)) * SCORE_KING_DEFENDERDEFICIT
			}
		}

		score = (score * factor[phase]) / 8
	} else {
		score += EndingKing[sq]
		score += ScoreControl(board, sq, side)

		pawns := board.Board[White][Pawn] | board.Board[Black][Pawn]
		for pawns != NullBitBoard {
			sq1 := LeadBit(pawns)
			ClearBit(&pawns, int(sq1))
			if BitPosArray[sq1]&board.Board[White][Pawn] != NullBitBoard {
				score -= Distance[sq][sq1+8]*10 - 5
			} else if BitPosArray[sq1]&board.Board[Black][Pawn] != NullBitBoard {
				score -= Distance[sq][sq1-8]*10 - 5
			} else {
				score -= Distance[sq][sq1] - 5
			}
		}
	}

	if phase >= 4 {
		if side == White {
			if sq < A2 {
				if MoveArray[King][sq]&(^board.Board[side][Pawn]&RankBit[1]) == NullBitBoard {
					score += SCORE_KING_BACK_RANK_WEAK
				}
			}
		} else {
			if sq > H7 {
				if MoveArray[King][sq]&(^board.Board[side][Pawn]&RankBit[6]) == NullBitBoard {
					score += SCORE_KING_BACK_RANK_WEAK
				}
			}
		}
	}

	return score
}

func BishopTrapped(board ChessBoard, side uint16) int {
	/* Check for trapped bishop at A2/H2/A7/H7 */
	score := 0
	if board.Board[side][Bishop] == NullBitBoard {
		return score
	}

	/* NOTE: SwapOff is a static exchange evaluator. It sees if the swap chain from the move is beneficial */
	if side == White {
		if (board.Board[White][Bishop]&BitPosArray[A7]) != 0 && (board.Board[Black][Pawn]&BitPosArray[B6]) != 0 && SwapOff(board, (A7<<6)|B6) < 0 {
			score += SCORE_BISHOP_TRAPPED
		}
		if (board.Board[White][Bishop]&BitPosArray[H7]) != 0 && (board.Board[Black][Pawn]&BitPosArray[G6]) != 0 && SwapOff(board, (H7<<6)|G6) < 0 {
			score += SCORE_BISHOP_TRAPPED
		}
	} else {
		if (board.Board[Black][Bishop]&BitPosArray[A2]) != 0 && (board.Board[White][Pawn]&BitPosArray[B3]) != 0 && SwapOff(board, (A2<<6)|B3) < 0 {
			score += SCORE_BISHOP_TRAPPED
		}
		if (board.Board[White][Bishop]&BitPosArray[H2]) != 0 && (board.Board[White][Pawn]&BitPosArray[G3]) != 0 && SwapOff(board, (H2<<6)|G3) < 0 {
			score += SCORE_BISHOP_TRAPPED
		}
	}
	return score
}

func DoubleQR7(board ChessBoard, side uint16) int {
	/* Check for QQ, QR, RR combo on 7th rank. */
	xside := 1 ^ side

	if NumBits((board.Board[side][Queen]|board.Board[side][Rook])&brank7[side]) > 1 && ((board.Board[xside][King]&brank8[side]) != 0 || (board.Board[xside][Pawn]&brank7[side]) != 0) {
		return SCORE_ROOK_7RANK
	}
	return 0
}

func ScorePiece(piece int, board ChessBoard, side uint16, pinnedBoard BitBoard) int {
	switch piece {
	case Knight:
		return ScoreKnight(board, side, pinnedBoard)
	case Bishop:
		return ScoreBishop(board, side, pinnedBoard)
	case Rook:
		return ScoreRook(board, side, pinnedBoard)
	case Queen:
		return ScoreQueen(board, side, pinnedBoard)
	}

	return 0
}

func ScoreControl(board ChessBoard, sq uint16, side uint16) int {
	/* Scoring for control side has on the board */
	score := 0

	enemyKing := board.KingPos[1^side]
	ourKing := board.KingPos[side]

	controlled := AttackXFrom(board, sq, side)

	n := NumBits(controlled & Boxes[0])
	score += 4 * int(n)

	n = NumBits(controlled & DistMap[enemyKing][2])
	score += int(n)

	n = NumBits(controlled & DistMap[ourKing][2])
	score += int(n)

	n = NumBits(controlled)
	score += 4 * int(n)

	return score
}
