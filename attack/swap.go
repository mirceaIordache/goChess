package attack

import (
	//	"fmt"

	. "github.com/mirceaIordache/goChess/common"
)

func SwapOff(board ChessBoard, move int) int {
	/* The SEE */
	from := MoveFrom(move)
	to := MoveTo(move)
	var side uint16
	var swaplist [MaxPlyDepth + 1]int
	var lastVal int
	cboard := GenerateCBoard(board)
	if board.Friends[White]&BitPosArray[from] != NullBitBoard {
		side = White
	} else {
		side = Black
	}
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
		swaplist[0] = int(Values[PromotePiece(move)])
		lastVal = -swaplist[0]
		swaplist[0] -= ValuePawn
	} else {
		if move&EnPassantMove != 0 {
			swaplist[0] = ValuePawn
		} else {
			swaplist[0] = int(Values[cboard[to]])
			lastVal = -int(Values[cboard[from]])
		}
	}
	plyDepth := 1
	for plyDepth < MaxPlyDepth {
		if attackingXSide == NullBitBoard {
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
				lastVal = int(Values[piece])
				break
			}
		}

		if attackingSide == NullBitBoard {
			break
		}
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
				lastVal = -int(Values[piece])
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

	return swaplist[0]
}

func AddXRayPiece(board ChessBoard, attackingSide, attackingXSide *BitBoard, to, sq, side uint16) {
	dir := Directions[to][sq]
	if dir == -1 {
		return
	}
	blockerRay := Ray[sq][dir] & board.Blocker
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
			*attackingXSide |= BitPosArray[numSq]
		}
	}
}
