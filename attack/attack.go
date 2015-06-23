package attack

import (
	. "github.com/mirceaIordache/goChess/common"
)

func AttackXFrom(board ChessBoard, sq uint16, side uint16) BitBoard {
	/* Attack pattern calculator for given piece */
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s, sq %d, side %s", ToEPD(board), sq, side)
	result := NullBitBoard

	pieces := board.Board[side]
	piece := GenerateCBoard(board)[sq]
	blocker := board.Blocker

	ChessLogger.Debug("Piece is %s", Notation[piece])

	switch piece {
	case Pawn:
		result = MoveArray[PawnType[side]][sq]

	case Knight:
		result = MoveArray[Knight][sq]

	case Bishop:
		fallthrough

	case Queen:
		newBlocker := blocker & (^(pieces[Bishop] | pieces[Queen]))
		for direction := RayBegin[Bishop]; direction < RayEnd[Bishop]; direction++ {
			moves := Ray[sq][direction] & newBlocker
			if moves == NullBitBoard {
				moves = Ray[sq][direction]
			} else {
				var blockedSquare uint16
				if BitPosArray[sq] > moves {
					blockedSquare = LeadBit(moves)
				} else {
					blockedSquare = TrailBit(moves)
				}

				moves = FromToRay[sq][blockedSquare]
			}
			result |= moves
		}

		fallthrough

	case Rook:
		if piece == Bishop {
			break
		}
		newBlocker := blocker & (^(pieces[Rook] | pieces[Queen]))
		for direction := RayBegin[Rook]; direction < RayEnd[Rook]; direction++ {
			moves := Ray[sq][direction] & newBlocker
			if moves == NullBitBoard {
				moves = Ray[sq][direction]
			} else {
				var blockedSquare uint16
				if BitPosArray[sq] > moves {
					blockedSquare = LeadBit(moves)
				} else {
					blockedSquare = TrailBit(moves)
				}
				moves = FromToRay[sq][blockedSquare]
			}
			result |= moves
		}

	case King:
		result = MoveArray[King][sq]
	}

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result is %064b", result)

	return result
}

func AttackTo(board ChessBoard, sq, side uint16) BitBoard {
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s, sq %d, side %d", ToEPD(board), sq, side)
	result := NullBitBoard
	attacker := board.Board[side]

	result |= (attacker[Knight] & MoveArray[Knight][sq])

	result |= (attacker[King] & MoveArray[King][sq])

	result |= (attacker[Pawn] & MoveArray[PawnType[1^side]][sq])

	ray := FromToRay[sq]
	blocker := board.Blocker

	throwAway := (attacker[Bishop] | attacker[Queen]) & MoveArray[Bishop][sq]

	for throwAway != NullBitBoard {
		sq1 := LeadBit(throwAway)
		ClearBit(&throwAway, int(sq1))
		if ray[sq1]&blocker&NotBitPosArray[sq1] == NullBitBoard {
			result |= BitPosArray[sq1]
		}
	}

	throwAway = (attacker[Rook] | attacker[Queen]) & MoveArray[Rook][sq]
	for throwAway != NullBitBoard {
		sq1 := LeadBit(throwAway)
		ClearBit(&throwAway, int(sq1))
		if ray[sq1]&blocker&NotBitPosArray[sq1] == NullBitBoard {
			result |= BitPosArray[sq1]
		}
	}

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result is %064b", result)
	return result
}

func GenerateAttacks(board ChessBoard) [2][7]BitBoard {
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s", ToEPD(board))

	/* Attacked pieces calculator */
	var result [2][7]BitBoard

	for side := White; side <= Black; side++ {
		attacked := board.Board[side]
		throwAway := &result[side][Knight]
		pieces := attacked[Knight]
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			*throwAway |= MoveArray[Knight][sq]
		}

		throwAway = &result[side][Bishop]
		pieces = attacked[Bishop]
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			*throwAway |= BishopAttack(board, sq)
		}

		throwAway = &result[side][Rook]
		pieces = attacked[Rook]
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			*throwAway |= RookAttack(board, sq)
		}

		throwAway = &result[side][Queen]
		pieces = attacked[Queen]
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			*throwAway |= QueenAttack(board, sq)
		}

		throwAway = &result[side][King]
		pieces = attacked[King]
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))
			*throwAway |= MoveArray[King][sq]
		}

		throwAway = &result[side][Pawn]
		if side == White {
			pieces = board.Board[White][Pawn] & ^FileBit[0]
			*throwAway |= pieces >> 7
			pieces = board.Board[White][Pawn] & ^FileBit[7]
			*throwAway |= pieces >> 9
		} else {
			pieces = board.Board[Black][Pawn] & ^FileBit[0]
			*throwAway |= pieces << 9
			pieces = board.Board[Black][Pawn] & ^FileBit[7]
			*throwAway |= pieces << 7
		}

		result[side][0] = result[side][Pawn] | result[side][Knight] | result[side][Bishop] | result[side][Rook] | result[side][Queen] | result[side][King]
	}

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result is %q", result[0:][0:])
	return result
}

func BishopAttack(board ChessBoard, sq uint16) BitBoard {
	ChessLogger.Info("Calculating Bishop Attacks")
	ChessLogger.Debug("Board %s, sq %d", ToEPD(board), sq)
	return Bishop45Attack[63-sq][(board.BlockerR45>>Shift45[63-sq])&Mask45[63-sq]] |
		Bishop315Attack[63-sq][(board.BlockerR315>>Shift315[63-sq])&Mask315[63-sq]]
}

func RookAttack(board ChessBoard, sq uint16) BitBoard {

	ChessLogger.Info("Calculating Rook Attacks")
	ChessLogger.Debug("Board %s, sq %d", ToEPD(board), sq)
	return Rook00Attack[63-sq][(board.Blocker>>Shift00[63-sq])&0xFF] |
		Rook90Attack[63-sq][(board.BlockerR90>>Shift90[63-sq])&0xFF]
}

func QueenAttack(board ChessBoard, sq uint16) BitBoard {

	ChessLogger.Info("Calculating Queen Attacks")
	ChessLogger.Debug("Board %s, sq %d", ToEPD(board), sq)
	return BishopAttack(board, sq) | RookAttack(board, sq)
}

func FindPins(board ChessBoard) BitBoard {
	ChessLogger.Info("Entering")
	/* Pinned Pieces calculator */
	result := NullBitBoard

	totalPieces := board.Friends[White] | board.Friends[Black]
	for side := White; side <= Black; side++ {
		xside := 1 ^ side
		attacks := GenerateAttacks(board)
		pins := board.Board[xside]

		attackees := pins[Rook] | pins[Queen] | pins[King]
		attackees |= (pins[Bishop] | pins[Knight]) & ^attacks[xside][0]

		pieces := board.Board[side][Bishop]
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))

			moves := MoveArray[Bishop][sq] & attackees
			for moves != NullBitBoard {
				sq1 := LeadBit(moves)
				ClearBit(&moves, int(sq1))
				attacker := totalPieces & NotBitPosArray[sq] & FromToRay[sq1][sq]
				if attacker&board.Friends[xside] != NullBitBoard && NumBits(attacker) == 1 {
					result |= attacker
				}
			}
		}

		attackees = pins[Queen] | pins[King]
		attackees |= (pins[Rook] | pins[Bishop] | pins[Knight]) & ^attacks[xside][0]

		pieces = board.Board[side][Rook]
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))

			moves := MoveArray[Rook][sq] & attackees
			for moves != NullBitBoard {
				sq1 := LeadBit(moves)
				ClearBit(&moves, int(sq1))
				attacker := totalPieces & NotBitPosArray[sq] & FromToRay[sq1][sq]
				if attacker&board.Friends[xside] != NullBitBoard && NumBits(attacker) == 1 {
					result |= attacker
				}
			}
		}

		attackees = pins[King]
		attackees |= (pins[Queen] | pins[Rook] | pins[Bishop] | pins[Knight]) & ^attacks[xside][0]

		pieces = board.Board[side][Queen]
		for pieces != NullBitBoard {
			sq := LeadBit(pieces)
			ClearBit(&pieces, int(sq))

			moves := MoveArray[Queen][sq] & attackees
			for moves != NullBitBoard {
				sq1 := LeadBit(moves)
				ClearBit(&moves, int(sq1))
				attacker := totalPieces & NotBitPosArray[sq] & FromToRay[sq1][sq]
				if attacker&board.Friends[xside] != NullBitBoard && NumBits(attacker) == 1 {
					result |= attacker
				}
			}
		}
	}
	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result is %064b", result)

	return result
}

func PinnedOnKing(board ChessBoard, sq, side uint16) bool {
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s, sq %d, side %d", ToEPD(board), sq, side)
	xside := 1 ^ side
	kingSq := board.KingPos[side]
	blocker := board.Blocker
	dir := Directions[kingSq][sq]

	if dir == -1 {
		ChessLogger.Info("Exiting")
		ChessLogger.Debug("Result %t", false)
		return false
	}

	if FromToRay[kingSq][sq]&NotBitPosArray[sq]&blocker != NullBitBoard {
		ChessLogger.Info("Exiting")
		ChessLogger.Debug("Result %t", false)

		return false
	}

	throwAway := (Ray[kingSq][dir] ^ FromToRay[kingSq][sq]) & blocker
	if throwAway == NullBitBoard {
		ChessLogger.Info("Exiting")
		ChessLogger.Debug("Result %t", false)

		return false
	}
	sq1 := uint16(0)
	if sq > kingSq {
		sq1 = LeadBit(throwAway)
	} else {
		sq1 = TrailBit(throwAway)
	}

	if dir <= 3 && BitPosArray[sq1]&(board.Board[xside][Queen]|board.Board[xside][Bishop]) != NullBitBoard {
		ChessLogger.Info("Exiting")
		ChessLogger.Debug("Result %t", true)

		return true
	}

	if dir >= 4 && BitPosArray[sq1]&(board.Board[xside][Queen]|board.Board[xside][Rook]) != NullBitBoard {
		ChessLogger.Info("Exiting")
		ChessLogger.Debug("Result %t", true)

		return true
	}

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result %t", false)

	return false
}

func SquareAttacked(board ChessBoard, sq, side uint16) bool {
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s, sq %d, side %d", ToEPD(board), sq, side)
	attacker := board.Board[side]
	if attacker[Knight]&MoveArray[Knight][sq] != NullBitBoard {
		ChessLogger.Info("Exiting")
		ChessLogger.Debug("Result %t", true)
		return true
	}

	if attacker[King]&MoveArray[King][sq] != NullBitBoard {
		ChessLogger.Info("Exiting")
		ChessLogger.Debug("Result %t", true)
		return true
	}

	if attacker[Pawn]&MoveArray[PawnType[1^side]][sq] != NullBitBoard {
		ChessLogger.Info("Exiting")
		ChessLogger.Debug("Result %t", true)
		return true
	}

	ray := FromToRay[sq]
	blocker := board.Blocker

	bishops := (attacker[Bishop] | attacker[Queen]) & MoveArray[Bishop][sq]
	deny := ^bishops & blocker

	for bishops != NullBitBoard {
		sq1 := LeadBit(bishops)
		if ray[sq1]&deny == NullBitBoard {
			ChessLogger.Info("Exiting")
			ChessLogger.Debug("Result %t", true)
			return true
		}
		ClearBit(&bishops, int(sq1))
	}

	rooks := (attacker[Rook] | attacker[Queen]) & MoveArray[Rook][sq]
	deny = ^rooks & blocker

	for rooks != NullBitBoard {
		sq1 := LeadBit(rooks)
		if ray[sq1]&deny == NullBitBoard {
			ChessLogger.Info("Exiting")
			ChessLogger.Debug("Result %t", true)
			return true
		}
		ClearBit(&rooks, int(sq1))
	}

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result %t", false)
	return false
}
