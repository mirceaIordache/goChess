package attack

import . "github.com/mirceaIordache/goChess/common"

func SwapOff(board ChessBoard, move int) int {
	/* The SEE */
	from := MoveFrom(move)
	to := MoveTo(move)
	ChessLogger.Info("Entering")
	var side uint16
	var swaplist [MaxPlyDepth]int
	var lastVal int
	cboard := GenerateCBoard(board)
	if board.Friends[White]&BitPosArray[from] != NullBitBoard {
		side = White
	} else {
		side = Black
	}
	ChessLogger.Debug("Board %s, From %s To %s, side %d", ToEPD(board), Algebraic[from], Algebraic[to], side)
	xside := 1 ^ side

	attackingSide := AttackTo(board, to, side)
	attackingXSide := AttackTo(board, to, xside)
	ClearBit(&attackingSide, int(from))
	if XRay[cboard[from]] != 0 {
		AddXRayPiece(board, &attackingSide, &attackingXSide, to, from, side)
	}

	boardSide := board.Board[side]
	boardXSide := board.Board[xside]

	if move&PiecePromotion != 0 {
		switch PromotePiece(move) {
		case 1:
			swaplist[0] = ValuePawn
		case 2:
			swaplist[0] = ValueKnight
		case 3:
			swaplist[0] = ValueBishop
		case 4:
			swaplist[0] = ValueRook
		case 5:
			swaplist[0] = ValueQueen
		}
		lastVal = -swaplist[0]
		swaplist[0] -= ValuePawn
	} else {
		if move&EnPassantMove != 0 {
			swaplist[0] = ValuePawn
		} else {
			switch cboard[to] {
			case Pawn:
				swaplist[0] = ValuePawn
			case Knight:
				swaplist[0] = ValueKnight
			case Bishop:
				swaplist[0] = ValueBishop
			case Rook:
				swaplist[0] = ValueRook
			case Queen:
				swaplist[0] = ValueQueen
			case King:
				swaplist[0] = ValueKing
			}
			switch cboard[from] {
			case Pawn:
				lastVal = -ValuePawn
			case Knight:
				lastVal = -ValueKnight
			case Bishop:
				lastVal = -ValueBishop
			case Rook:
				lastVal = -ValueRook
			case Queen:
				lastVal = -ValueQueen
			case King:
				lastVal = -ValueKing
			}
		}
	}
	plyDepth := 1
	for plyDepth < MaxPlyDepth/2-1 {
		ChessLogger.Debug("Iteration %d, side %d, bitboard %064b", plyDepth, xside, attackingXSide)
		if attackingXSide == NullBitBoard {
			ChessLogger.Info("No counter move. Exiting")
			break
		}
		for piece := Pawn; piece <= King; piece++ {
			xsidePieces := attackingXSide & boardXSide[piece]
			if xsidePieces != NullBitBoard {
				sq := LeadBit(xsidePieces)
				ClearBit(&xsidePieces, int(sq))
				if XRay[piece] != 0 {
					AddXRayPiece(board, &attackingXSide, &attackingSide, to, sq, xside)
				}
				swaplist[plyDepth] = swaplist[plyDepth-1] + lastVal
				plyDepth++
				switch piece {
				case Pawn:
					lastVal = ValuePawn
				case Knight:
					lastVal = ValueKnight
				case Bishop:
					lastVal = ValueBishop
				case Rook:
					lastVal = ValueRook
				case Queen:
					lastVal = ValueQueen
				case King:
					lastVal = ValueKing
				}
				break
			}
		}

		if attackingSide == NullBitBoard {
			ChessLogger.Info("No counter move. Exiting")
			break
		}
		ChessLogger.Debug("Iteration %d, side %d", plyDepth, side)

		for piece := Pawn; piece <= King; piece++ {
			sidePieces := attackingSide & boardSide[piece]
			if sidePieces != NullBitBoard {
				sq := LeadBit(sidePieces)
				ClearBit(&sidePieces, int(sq))
				if XRay[piece] != 0 {
					AddXRayPiece(board, &attackingSide, &attackingXSide, to, sq, side)
				}
				swaplist[plyDepth] = swaplist[plyDepth-1] + lastVal
				plyDepth++
				switch piece {
				case Pawn:
					lastVal = -ValuePawn
				case Knight:
					lastVal = -ValueKnight
				case Bishop:
					lastVal = -ValueBishop
				case Rook:
					lastVal = -ValueRook
				case Queen:
					lastVal = -ValueQueen
				case King:
					lastVal = -ValueKing
				}
				break
			}
		}
	}

	for plyDepth--; plyDepth != 0; plyDepth-- {
		if plyDepth&1 != 0 {
			if swaplist[plyDepth] <= swaplist[plyDepth-1] {
				swaplist[plyDepth-1] = swaplist[plyDepth]
			}
		} else {
			if swaplist[plyDepth] >= swaplist[plyDepth-1] {
				swaplist[plyDepth-1] = swaplist[plyDepth]
			}
		}
	}

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result %d", swaplist[0])
	return swaplist[0]
}

func AddXRayPiece(board ChessBoard, attackingSide, attackingXSide *BitBoard, to, sq, side uint16) {
	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s, AttackingSide %064b, AttackingXSide %064b, from %s to %s, side %d", ToEPD(board), *attackingSide, *attackingXSide, Algebraic[sq], Algebraic[to], side)
	dir := Directions[to][sq]
	blockerRay := Ray[63-sq][dir] & board.Blocker
	if blockerRay == NullBitBoard {
		return
	}

	var numSq uint16
	if to < sq {
		numSq = LeadBit(blockerRay)
	} else {
		numSq = TrailBit(blockerRay)
	}
	piece := GenerateCBoard(board)[numSq]
	
	if piece == Queen || (piece == Rook && dir > 3) || (piece == Bishop && dir < 4) {
		if BitPosArray[numSq]&board.Friends[side] != NullBitBoard {
			*attackingSide |= BitPosArray[numSq]
		} else {
			*attackingXSide &= ^BitPosArray[numSq]
		}
	}

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result attackingSide %064b, attackingXSide %064b", *attackingSide, *attackingXSide)
}
