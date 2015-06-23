package common

import (
	"math"
)

func GenMove(from, to uint16) int {
	return int((from << 6) | to)
}

func MoveFrom(move int) uint16 {
	return uint16((move >> 6) & 0x003F)
}

func MoveTo(move int) uint16 {
	return uint16(move & 0x003F)
}

func PromotePiece(piece int) int {
	return (piece >> 12) & 0x0007
}

func CapturePiece(piece int) int {
	return (piece >> 15) & 0x0007
}

func ApplyNullMove(board ChessBoard, move int, side uint16) ChessBoard {
	return MakeMove(board, move, side, true)
}

func ApplyMove(board ChessBoard, move int, side uint16) ChessBoard {
	return MakeMove(board, move, side, false)
}

func MakeMove(board ChessBoard, move int, side uint16, null bool) ChessBoard {
	//Application of move on chess board

	newBoard := board
	xside := 1 ^ side
	from := MoveFrom(move)
	to := MoveTo(move)
	cboard := GenerateCBoard(newBoard)
	fromPiece := cboard[from]
	toPiece := cboard[to]

	ChessLogger.Info("Entering")
	ChessLogger.Debug("Board %s, From %s to %s, side %d, null move %t", ToEPD(board), Algebraic[MoveFrom(move)], Algebraic[MoveTo(move)], side, null)

	moveBitBoard := &newBoard.Board[side][fromPiece]
	ClearBit(moveBitBoard, int(from))
	SetBit(moveBitBoard, int(to))
	ClearBit(&newBoard.BlockerR90, R90[from])
	SetBit(&newBoard.BlockerR90, R90[to])
	ClearBit(&newBoard.BlockerR45, R45[from])
	SetBit(&newBoard.BlockerR45, R45[to])
	ClearBit(&newBoard.BlockerR315, R315[from])
	SetBit(&newBoard.BlockerR315, R315[to])

	if fromPiece == King {
		board.KingPos[side] = to
	}

	if toPiece != Empty {
		// Capture
		ClearBit(&newBoard.Board[xside][toPiece], int(to))
		newBoard.Material[xside] -= Values[toPiece]
		if toPiece != Pawn {
			newBoard.PawnMaterial[xside] -= Values[toPiece]
		}
	}

	if move&PiecePromotion != 0 {
		// Promotion
		SetBit(&newBoard.Board[side][PromotePiece(move)], int(to))
		ClearBit(moveBitBoard, int(to))
		newBoard.Material[side] += Values[cboard[to]] - ValuePawn
		newBoard.PawnMaterial[side] += Values[cboard[to]]
	}

	if move&EnPassantMove != 0 {
		//En passant
		delta := 0
		if side == White {
			delta = -8
		} else {
			delta = 8
		}
		epSq := board.Ep + int16(delta)
		ClearBit(&newBoard.Board[xside][Pawn], int(epSq))
		ClearBit(&newBoard.BlockerR90, R90[epSq])
		ClearBit(&newBoard.BlockerR45, R45[epSq])
		ClearBit(&newBoard.BlockerR315, R315[epSq])
		newBoard.Material[xside] -= ValuePawn
	}

	if move&CastlingMove != 0 {
		rookFrom := uint16(0)
		rookTo := uint16(0)
		if to&0x04 != 0 {
			rookFrom = to + 1
			rookTo = to - 1
		} else {
			rookFrom = to - 2
			rookTo = to + 1
		}
		ClearBit(&newBoard.Board[side][Rook], int(rookFrom))
		SetBit(&newBoard.Board[side][Rook], int(rookTo))
		ClearBit(&newBoard.BlockerR90, R90[rookFrom])
		SetBit(&newBoard.BlockerR90, R90[rookTo])
		ClearBit(&newBoard.BlockerR45, R45[rookFrom])
		SetBit(&newBoard.BlockerR45, R45[rookTo])
		ClearBit(&newBoard.BlockerR315, R315[rookFrom])
		SetBit(&newBoard.BlockerR315, R315[rookTo])
		newBoard.Castled[side] = true
	}

	if side == White {
		if fromPiece == King && newBoard.CastleFlag&WCastle != 0 {
			newBoard.CastleFlag &= ^WCastle
		} else if fromPiece == Rook && null == false {
			if from == H1 {
				newBoard.CastleFlag &= ^WKingCastle
			} else if from == A1 {
				newBoard.CastleFlag &= ^WQueenCastle
			}
		}
		if toPiece == Rook && null == false {
			if to == H8 {
				newBoard.CastleFlag &= ^BKingCastle
			} else if to == A8 {
				newBoard.CastleFlag &= ^BQueenCastle
			}
		}
	} else {
		if fromPiece == King && newBoard.CastleFlag&BCastle != 0 {
			newBoard.CastleFlag &= ^BCastle
		} else if fromPiece == Rook && null == false {
			if from == H8 {
				newBoard.CastleFlag &= ^BKingCastle
			} else if from == A8 {
				newBoard.CastleFlag &= ^BQueenCastle
			}
		}
		if toPiece == Rook && null == false {
			if to == H1 {
				newBoard.CastleFlag &= ^WKingCastle
			} else if to == A8 {
				newBoard.CastleFlag &= ^WQueenCastle
			}
		}
	}

	if fromPiece == Pawn && math.Abs(float64(from-to)) == float64(16) {
		newBoard.Ep = int16((from + to) / 2)
	} else {
		newBoard.Ep = -1
	}

	newBoard.Side = xside

	newBoard.Friends[White] = newBoard.Board[White][Pawn] | newBoard.Board[White][Knight] | newBoard.Board[White][Bishop] | newBoard.Board[White][Rook] | newBoard.Board[White][Queen] | newBoard.Board[White][King]
	newBoard.Friends[Black] = newBoard.Board[Black][Pawn] | newBoard.Board[Black][Knight] | newBoard.Board[Black][Bishop] | newBoard.Board[Black][Rook] | newBoard.Board[Black][Queen] | newBoard.Board[Black][King]
	newBoard.Blocker = newBoard.Friends[White] | newBoard.Friends[Black]

	ChessLogger.Info("Exiting")
	ChessLogger.Debug("Result %s", ToEPD(newBoard))

	return newBoard
}
