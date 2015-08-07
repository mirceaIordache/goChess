package moveGenerator

import (
	. "github.com/mirceaIordache/goChess/attack"
	. "github.com/mirceaIordache/goChess/common"
)

type MoveList struct {
	Next  *MoveList
	Value Move
}

func AddMoveToList(list *MoveList, from, to uint16, mask int) {
	list.Value.Move = GenMove(from, to) | mask
	list.Next = &MoveList{nil, Move{}}
	list = list.Next
}

func AddPromotionToList(list *MoveList, from, to uint16) {
	AddMoveToList(list, from, to, QueenPromotion)
	AddMoveToList(list, from, to, KnightPromotion)
	AddMoveToList(list, from, to, RookPromotion)
	AddMoveToList(list, from, to, BishopPromotion)
}

func BitToMove(board BitBoard, sq uint16, list *MoveList) {

	for board != NullBitBoard {
		sq1 := LeadBit(board)
		ClearBit(&board, int(sq1))
		AddMoveToList(list, sq, sq1, 0)
	}
}

func GenerateMoves(board ChessBoard) *MoveList {
	/* Generate all possible moves */
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s", ToEPD(board))
	var tail *MoveList = &MoveList{nil, Move{}}
	head := tail

	side := board.Side
	ourPieces := board.Board[side]
	friends := board.Friends[side]
	notFriends := 1 ^ friends
	blocker := board.Blocker
	notBlocker := ^blocker
	enPassant := board.Ep

	/* Knight & King */
	for piece := Knight; piece <= King; piece += 4 {
		pieces := ourPieces[piece]
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			BitToMove(MoveArray[piece][sq]&notFriends, sq, tail)
		}
	}

	/* Bishop */
	pieces := ourPieces[Bishop]
	for pieces != NullBitBoard {
		sq := LeadBit(pieces)
		ClearBit(&pieces, int(sq))
		pattern := BishopAttack(board, sq)
		BitToMove(pattern&notFriends, sq, tail)
	}

	/* Rook */
	pieces = ourPieces[Rook]
	for pieces != NullBitBoard {
		sq := LeadBit(pieces)
		ClearBit(&pieces, int(sq))
		pattern := RookAttack(board, sq)
		BitToMove(pattern&notFriends, sq, tail)
	}

	/* Queen */
	pieces = ourPieces[Queen]
	for pieces != NullBitBoard {
		sq := LeadBit(pieces)
		ClearBit(&pieces, int(sq))
		pattern := QueenAttack(board, sq)
		BitToMove(pattern&notFriends, sq, tail)
	}

	var epArray BitBoard
	if enPassant > -1 {
		epArray = BitPosArray[enPassant]
	} else {
		epArray = NullBitBoard
	}

	/* Pawns */
	enemies := board.Friends[1^side] | epArray
	if side == White {
		/* 1 square forward */
		throwAway := (ourPieces[Pawn] << 8) & notBlocker
		for throwAway != NullBitBoard {
			sq := LeadBit(throwAway)
			ClearBit(&throwAway, int(sq))
			if sq >= 56 {
				AddPromotionToList(tail, sq-8, sq)
			} else {
				AddMoveToList(tail, sq-8, sq, 0)
			}
		}

		/* Pawns on 2nd rank 2 squares forward */
		pieces = ourPieces[Pawn] & RankBit[1]
		throwAway = (pieces << 8) & notBlocker
		throwAway = (throwAway << 8) & notBlocker
		for throwAway != NullBitBoard {
			sq := LeadBit(throwAway)
			ClearBit(&throwAway, int(sq))
			AddMoveToList(tail, sq-16, sq, 0)
		}

		/* Attacks to the left */
		pieces = ourPieces[Pawn] & ^FileBit[0]
		throwAway = (pieces << 7) & enemies
		for throwAway != NullBitBoard {
			sq := LeadBit(throwAway)
			ClearBit(&throwAway, int(sq))
			if sq >= 56 {
				AddPromotionToList(tail, sq-7, sq)
			} else if uint16(enPassant) == sq {
				AddMoveToList(tail, sq-7, sq, EnPassantMove)
			} else {
				AddMoveToList(tail, sq-7, sq, 0)
			}
		}

		/* Attacks to the right */
		pieces = ourPieces[Pawn] & ^FileBit[7]
		throwAway = (pieces << 9) & enemies
		for throwAway != NullBitBoard {
			sq := LeadBit(throwAway)
			ClearBit(&throwAway, int(sq))
			if sq >= 56 {
				AddPromotionToList(tail, sq-9, sq)
			} else if uint16(enPassant) == sq {
				AddMoveToList(tail, sq-9, sq, EnPassantMove)
			} else {
				AddMoveToList(tail, sq-9, sq, 0)
			}
		}
	}

	/* Same as for White side */
	if side == Black {
		throwAway := (ourPieces[Pawn] >> 8) & notBlocker
		for throwAway != NullBitBoard {
			sq := LeadBit(throwAway)
			ClearBit(&throwAway, int(sq))
			if sq <= 7 {
				AddPromotionToList(tail, sq+8, sq)
			} else {
				AddMoveToList(tail, sq+8, sq, 0)
			}
		}

		pieces = ourPieces[Pawn] & RankBit[6]
		throwAway = (pieces >> 8) & notBlocker
		throwAway = (throwAway >> 8) & notBlocker
		for throwAway != NullBitBoard {
			sq := LeadBit(throwAway)
			ClearBit(&throwAway, int(sq))
			AddMoveToList(tail, sq+16, sq, 0)
		}

		pieces = ourPieces[Pawn] & ^FileBit[7]
		throwAway = (pieces >> 7) & enemies
		for throwAway != NullBitBoard {
			sq := LeadBit(throwAway)
			ClearBit(&throwAway, int(sq))
			if sq <= 7 {
				AddPromotionToList(tail, sq+7, sq)
			} else if uint16(enPassant) == sq {
				AddMoveToList(tail, sq+7, sq, EnPassantMove)
			} else {
				AddMoveToList(tail, sq+7, sq, 0)
			}
		}

		pieces = ourPieces[Pawn] & ^FileBit[0]
		throwAway = (pieces >> 9) & enemies
		for throwAway != NullBitBoard {
			sq := LeadBit(throwAway)
			ClearBit(&throwAway, int(sq))
			if sq <= 7 {
				AddPromotionToList(tail, sq+9, sq)
			} else if uint16(enPassant) == sq {
				AddMoveToList(tail, sq+9, sq, EnPassantMove)
			} else {
				AddMoveToList(tail, sq+9, sq, 0)
			}
		}
	}

	/* Castling, if not done and possible */
	pieces = ourPieces[Rook]
	if side == White && (board.CastleFlag&WKingCastle != 0) && (pieces&BitPosArray[H1] != NullBitBoard) && (FromToRay[E1][G1]&blocker == NullBitBoard) && SquareAttacked(board, E1, Black) == false && SquareAttacked(board, F1, Black) == false && SquareAttacked(board, G1, Black) == false {
		AddMoveToList(tail, E1, G1, CastlingMove)
	}

	if side == White && (board.CastleFlag&WQueenCastle != 0) && (pieces&BitPosArray[A1] != NullBitBoard) && (FromToRay[E1][B1]&blocker == NullBitBoard) && SquareAttacked(board, E1, Black) == false && SquareAttacked(board, D1, Black) == false && SquareAttacked(board, C1, Black) == false {
		AddMoveToList(tail, E1, C1, CastlingMove)
	}

	if side == Black && (board.CastleFlag&BKingCastle != 0) && (pieces&BitPosArray[H8] != NullBitBoard) && (FromToRay[E8][G8]&blocker == NullBitBoard) && SquareAttacked(board, E8, White) == false && SquareAttacked(board, F8, White) == false && SquareAttacked(board, G8, White) == false {
		AddMoveToList(tail, E8, G8, CastlingMove)
	}

	if side == Black && (board.CastleFlag&BQueenCastle != 0) && (pieces&BitPosArray[A8] != NullBitBoard) && (FromToRay[E8][B8]&blocker == NullBitBoard) && SquareAttacked(board, E8, White) == false && SquareAttacked(board, D8, White) == false && SquareAttacked(board, C8, White) == false {
		AddMoveToList(tail, E8, C8, CastlingMove)
	}

	tail = head
	for tail.Next != nil {
		tail = tail.Next
	}

	ChessLogger.Info("Exiting")
	return head
}

func GenerateNonCaptures(board ChessBoard) *MoveList {
	/* Generate non capture moves */
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s", ToEPD(board))
	var tail *MoveList = &MoveList{nil, Move{}}
	head := tail

	side := board.Side
	ourPieces := board.Board[side]
	blocker := board.Blocker
	notBlocker := 1 ^ blocker

	/* Knight and King */
	for piece := Knight; piece <= King; piece += 4 {
		pieces := ourPieces[piece]
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			BitToMove(MoveArray[piece][sq]&notBlocker, sq, tail)
		}
	}

	/* Bishop */
	pieces := ourPieces[Bishop]
	for pieces != NullBitBoard {
		sq := LeadBit(pieces)
		ClearBit(&pieces, int(sq))
		pattern := BishopAttack(board, sq)
		BitToMove(pattern&notBlocker, sq, tail)
	}

	/* Rook */
	pieces = ourPieces[Rook]
	for pieces != NullBitBoard {
		sq := LeadBit(pieces)
		ClearBit(&pieces, int(sq))
		pattern := RookAttack(board, sq)
		BitToMove(pattern&notBlocker, sq, tail)
	}

	/* Queen */
	pieces = ourPieces[Queen]
	for pieces != NullBitBoard {
		sq := LeadBit(pieces)
		ClearBit(&pieces, int(sq))
		pattern := QueenAttack(board, sq)
		BitToMove(pattern&notBlocker, sq, tail)
	}

	/* White Pawn moves */
	if side == White {
		pieces := (ourPieces[Pawn] << 8) & notBlocker
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			if sq < 56 {
				AddMoveToList(tail, sq-8, sq, 0)
			}
		}

		pieces = ourPieces[Pawn] & RankBit[1]
		pieces = (pieces << 8) & notBlocker
		pieces = (pieces << 8) & notBlocker
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			AddMoveToList(tail, sq-16, sq, 0)

		}
	}

	if side == Black {
		pieces := (ourPieces[Pawn] >> 8) & notBlocker
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			if sq > 7 {
				AddMoveToList(tail, sq+8, sq, 0)
			}
		}

		pieces = ourPieces[Pawn] & RankBit[6]
		pieces = (pieces >> 8) & notBlocker
		pieces = (pieces >> 8) & notBlocker
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			AddMoveToList(tail, sq+16, sq, 0)

		}
	}

	/* Castling */
	pieces = ourPieces[Rook]
	if side == White && (board.CastleFlag&WKingCastle != 0) && (pieces&BitPosArray[H1] != NullBitBoard) && (FromToRay[E1][G1]&blocker == NullBitBoard) && SquareAttacked(board, E1, Black) == false && SquareAttacked(board, F1, Black) == false && SquareAttacked(board, G1, Black) == false {
		AddMoveToList(tail, E1, G1, CastlingMove)
	}

	if side == White && (board.CastleFlag&WQueenCastle != 0) && (pieces&BitPosArray[A1] != NullBitBoard) && (FromToRay[E1][B1]&blocker == NullBitBoard) && SquareAttacked(board, E1, Black) == false && SquareAttacked(board, D1, Black) == false && SquareAttacked(board, C1, Black) == false {
		AddMoveToList(tail, E1, C1, CastlingMove)
	}

	if side == Black && (board.CastleFlag&BKingCastle != 0) && (pieces&BitPosArray[H8] != NullBitBoard) && (FromToRay[E8][G8]&blocker == NullBitBoard) && SquareAttacked(board, E8, White) == false && SquareAttacked(board, F8, White) == false && SquareAttacked(board, G8, White) == false {
		AddMoveToList(tail, E8, G8, CastlingMove)
	}

	if side == Black && (board.CastleFlag&BQueenCastle != 0) && (pieces&BitPosArray[A8] != NullBitBoard) && (FromToRay[E8][B8]&blocker == NullBitBoard) && SquareAttacked(board, E8, White) == false && SquareAttacked(board, D8, White) == false && SquareAttacked(board, C8, White) == false {
		AddMoveToList(tail, E8, C8, CastlingMove)
	}

	ChessLogger.Info("Exiting")
	return head
}

func GenerateCaptures(board ChessBoard) *MoveList {
	/* Generate captures. Include en passant and promotions */
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s", ToEPD(board))
	var tail *MoveList = &MoveList{nil, Move{}}
	head := tail

	side := board.Side
	ourPieces := board.Board[side]
	enemies := board.Friends[1^side]
	blocker := board.Blocker
	enPassant := board.Ep

	/* Knight and King */
	for piece := Knight; piece <= King; piece += 4 {
		pieces := ourPieces[piece]
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			BitToMove(MoveArray[piece][sq]&enemies, sq, tail)
		}
	}

	/* Bishop */
	pieces := ourPieces[Bishop]
	for pieces != NullBitBoard {
		sq := LeadBit(pieces)
		ClearBit(&pieces, int(sq))
		pattern := BishopAttack(board, sq)
		BitToMove(pattern&enemies, sq, tail)
	}

	/* Rook */
	pieces = ourPieces[Rook]
	for pieces != NullBitBoard {
		sq := LeadBit(pieces)
		ClearBit(&pieces, int(sq))
		pattern := RookAttack(board, sq)
		BitToMove(pattern&enemies, sq, tail)
	}

	/* Queen */
	pieces = ourPieces[Queen]
	for pieces != NullBitBoard {
		sq := LeadBit(pieces)
		ClearBit(&pieces, int(sq))
		pattern := QueenAttack(board, sq)
		BitToMove(pattern&enemies, sq, tail)
	}

	var epArray BitBoard
	if enPassant > -1 {
		epArray = BitPosArray[enPassant]
	} else {
		epArray = NullBitBoard
	}

	/* White pawns */
	if side == White {
		/* Promotions */
		pieces = ourPieces[Pawn] & RankBit[6]
		pieces = (pieces << 8) & ^blocker
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			AddPromotionToList(tail, sq-8, sq)
		}

		/* Captures to the left */
		pieces = ourPieces[Pawn] & ^FileBit[0]
		pieces = (pieces << 7) & (board.Friends[1^side] | epArray)
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			if sq >= 56 {
				AddPromotionToList(tail, sq-7, sq)
			} else if enPassant == int16(sq) {
				AddMoveToList(tail, sq-7, sq, EnPassantMove)
			} else {
				AddMoveToList(tail, sq-7, sq, 0)
			}
		}

		pieces = ourPieces[Pawn] & ^FileBit[7]
		pieces = (pieces << 9) & (board.Friends[1^side] | epArray)
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			if sq >= 56 {
				AddPromotionToList(tail, sq-9, sq)
			} else if enPassant == int16(sq) {
				AddMoveToList(tail, sq-9, sq, EnPassantMove)
			} else {
				AddMoveToList(tail, sq-9, sq, 0)
			}
		}
	}

	/* Black Pawns */
	if side == Black {
		/* Promotions */
		pieces = ourPieces[Pawn] & RankBit[1]
		pieces = (pieces >> 8) & ^blocker
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			AddPromotionToList(tail, sq+8, sq)
		}

		/* Captures to the left */
		pieces = ourPieces[Pawn] & ^FileBit[7]
		pieces = (pieces >> 7) & (board.Friends[1^side] | epArray)
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			if sq <= 7 {
				AddPromotionToList(tail, sq+7, sq)
			} else if enPassant == int16(sq) {
				AddMoveToList(tail, sq+7, sq, EnPassantMove)
			} else {
				AddMoveToList(tail, sq+7, sq, 0)
			}
		}

		pieces = ourPieces[Pawn] & ^FileBit[0]
		pieces = (pieces >> 9) & (board.Friends[1^side] | epArray)
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			if sq <= 7 {
				AddPromotionToList(tail, sq+9, sq)
			} else if enPassant == int16(sq) {
				AddMoveToList(tail, sq+9, sq, EnPassantMove)
			} else {
				AddMoveToList(tail, sq+9, sq, 0)
			}
		}
	}

	ChessLogger.Info("Exiting")
	return head
}

func GenerateCheckEscapes(board ChessBoard) *MoveList {
	/* Generate moves that would get the king out of check */
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s", ToEPD(board))
	var head *MoveList = &MoveList{nil, Move{}}
	tail := head

	side := board.Side
	xside := 1 ^ side
	kingSq := board.KingPos[side]
	checkers := AttackTo(board, kingSq, xside)
	ourPawns := board.Board[side][Pawn]
	cboard := GenerateCBoard(board)

	if NumBits(checkers) == 1 {
		/* Try to capture checking piece */
		checkSq := LeadBit(checkers)
		attackers := AttackTo(board, checkSq, side)
		attackers &= ^board.Board[side][King]
		for attackers != NullBitBoard {
			sq := LeadBit(attackers)
			ClearBit(&attackers, int(sq))
			if PinnedOnKing(board, sq, side) == false {
				if GenerateCBoard(board)[sq] == Pawn && (checkSq <= H1 || checkSq >= A8) {
					AddPromotionToList(tail, sq, checkSq)
				} else {
					AddMoveToList(tail, sq, checkSq, 0)
				}
			}
		}

		/* Check if en passant can help */
		if board.Ep > -1 {
			epSq := board.Ep
			var sideDelta int16
			if side == White {
				sideDelta = -8
			} else {
				sideDelta = 8
			}

			if epSq+sideDelta == int16(checkSq) {
				throwAway := MoveArray[PawnType[1^side]][epSq] & ourPawns
				for throwAway != NullBitBoard {
					sq := LeadBit(throwAway)
					ClearBit(&throwAway, int(sq))
					if PinnedOnKing(board, sq, side) == false {
						AddMoveToList(tail, sq, uint16(epSq), EnPassantMove)
					}
				}
			}
		}

		/* Block or capture checking piece */
		if Slider[cboard[checkSq]] == 1 {
			throwAway := FromToRay[63-kingSq][63-checkSq] & NotBitPosArray[checkSq]
			for throwAway != NullBitBoard {
				sq := LeadBit(throwAway)
				ClearBit(&throwAway, int(sq))
				attackers := AttackTo(board, sq, side)
				attackers &= ^(board.Board[side][King] | ourPawns)

				/* Add in pawn advances */
				if side == White && sq > H2 {
					if BitPosArray[sq-8]&ourPawns != NullBitBoard {
						attackers |= BitPosArray[sq-8]
					}
					if sq>>3 == 3 && cboard[sq-8] == Empty && BitPosArray[sq-16]&ourPawns != NullBitBoard {
						attackers |= BitPosArray[sq-16]
					}
				}

				if side == Black && sq < A7 {
					if BitPosArray[sq+8]&ourPawns != NullBitBoard {
						attackers |= BitPosArray[sq+8]
					}

					if sq>>3 == 4 && cboard[sq+8] == Empty && BitPosArray[sq+16]&ourPawns != NullBitBoard {
						attackers |= BitPosArray[sq+16]
					}
				}

				for attackers != NullBitBoard {
					sq1 := LeadBit(attackers)
					ClearBit(&attackers, int(sq1))
					if PinnedOnKing(board, sq1, side) == false {
						if cboard[sq1] == Pawn && (sq > H7 || sq < A2) {
							AddPromotionToList(tail, sq1, sq)
						} else {
							AddMoveToList(tail, sq1, sq, 0)
						}
					}
				}
			}
		}
	}

	escapes := NullBitBoard
	if checkers != NullBitBoard {
		escapes = MoveArray[King][kingSq] & ^board.Friends[side]
	}

	for checkers != NullBitBoard {
		checkSq := LeadBit(checkers)
		ClearBit(&checkers, int(checkSq))
		direction := Directions[checkSq][kingSq]
		if Slider[cboard[checkSq]] == 1 {
			escapes &= ^Ray[checkSq][direction]
		}
	}

	for escapes != NullBitBoard {
		sq := LeadBit(escapes)
		ClearBit(&escapes, int(sq))
		if SquareAttacked(board, sq, xside) == false {
			AddMoveToList(tail, kingSq, sq, 0)
		}
	}

	ChessLogger.Info("Exiting")
	return head
}

func FilterMoves(board ChessBoard, list *MoveList) {
	/* Remove illegal moves from move list */
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s", ToEPD(board))
	var prevInList *MoveList = nil
	var check bool
	side := board.Side
	xside := 1 ^ side
	iter := list

	kingSq := board.KingPos[side]
	for iter.Next != nil {
		newBoard := ApplyMove(board, iter.Value.Move, side)
		if GenerateCBoard(newBoard)[MoveTo(iter.Value.Move)] != King {
			check = SquareAttacked(newBoard, kingSq, xside)
		} else {
			check = SquareAttacked(newBoard, MoveTo(iter.Value.Move), xside)
		}
		if check == true {
			*iter = *iter.Next
			if prevInList != nil {
				prevInList.Next = iter
			}
		}
		if iter.Next != nil {
			prevInList = iter
			iter = iter.Next
		}
	}
	ChessLogger.Info("Exiting")
}
