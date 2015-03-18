package common

import (
	"bytes"
	"fmt"
)

/** Used in serializing data to send over the network */
func ToEPD(board ChessBoard) string {
	//	fmt.Printf("Entered toEPD\n")
	var r, c, sq, k int
	var buffer bytes.Buffer

	//	fmt.Printf("Generating CBoard\n")
	cboard := GenerateCBoard(board)

	for r = A8; r >= A1; r -= 8 {
		k = 0
		for c = 0; c < 8; c++ {
			sq = r + c
			//fmt.Printf("SQ: %d\n", sq)

			if cboard[sq] == Empty {
				k++
			} else {
				if k != 0 {
					buffer.WriteString(fmt.Sprintf("%1d", k))
				}
				k = 0
				c1 := Notation[cboard[sq]]
				if (BitPosArray[sq] & board.Friends[Black]) != 0 {
					c1 = LNotation[cboard[sq]]
				}
				buffer.WriteString(fmt.Sprintf("%c", c1))
			}
		}
		if k != 0 {
			buffer.WriteString(fmt.Sprintf("%1d", k))
		}
		if r > A1 {
			buffer.WriteString("/")
		}

	}

	/* Print aux stuff */

	side := " w "
	if board.Side == Black {
		side = " b "
	}

	buffer.WriteString(side)

	if board.CastleFlag&WKingCastle != 0 {
		buffer.WriteString("K")
	}
	if board.CastleFlag&WQueenCastle != 0 {
		buffer.WriteString("Q")
	}
	if board.CastleFlag&BKingCastle != 0 {
		buffer.WriteString("k")
	}
	if board.CastleFlag&BQueenCastle != 0 {
		buffer.WriteString("q")
	}
	if board.CastleFlag&(WCastle|BCastle) == 0 {
		buffer.WriteString("-")
	}

	ep := " -"
	if board.Ep > -1 {
		ep = fmt.Sprintf(" %s", Algebraic[board.Ep])
	}
	buffer.WriteString(ep)

	return buffer.String()
}

/** Used in deserializing the data sent over network */
func FromEPD(boardString string) ChessBoard {
	board := new(ChessBoard)

	r := 56
	c := 0
	sq := 0
	var i int
	for j, char := range boardString {
		i = j
		if char == ' ' {
			break
		}
		sq = r + c
		switch char {
		case 'P':
			SetBit(&board.Board[White][Pawn], sq)
			SetBit(&board.BlockerR90, R90[sq])
			SetBit(&board.BlockerR45, R45[sq])
			SetBit(&board.BlockerR315, R315[sq])
			board.Material[White] += ValuePawn
		case 'N':
			SetBit(&board.Board[White][Knight], sq)
			SetBit(&board.BlockerR90, R90[sq])
			SetBit(&board.BlockerR45, R45[sq])
			SetBit(&board.BlockerR315, R315[sq])
			board.Material[White] += ValueKnight
		case 'B':
			SetBit(&board.Board[White][Bishop], sq)
			SetBit(&board.BlockerR90, R90[sq])
			SetBit(&board.BlockerR45, R45[sq])
			SetBit(&board.BlockerR315, R315[sq])
			board.Material[White] += ValueBishop
		case 'R':
			SetBit(&board.Board[White][Rook], sq)
			SetBit(&board.BlockerR90, R90[sq])
			SetBit(&board.BlockerR45, R45[sq])
			SetBit(&board.BlockerR315, R315[sq])
			board.Material[White] += ValueRook
		case 'Q':
			SetBit(&board.Board[White][Queen], sq)
			SetBit(&board.BlockerR90, R90[sq])
			SetBit(&board.BlockerR45, R45[sq])
			SetBit(&board.BlockerR315, R315[sq])
			board.Material[White] += ValueQueen
		case 'K':
			SetBit(&board.Board[White][King], sq)
			SetBit(&board.BlockerR90, R90[sq])
			SetBit(&board.BlockerR45, R45[sq])
			SetBit(&board.BlockerR315, R315[sq])
		case 'p':
			SetBit(&board.Board[Black][Pawn], sq)
			SetBit(&board.BlockerR90, R90[sq])
			SetBit(&board.BlockerR45, R45[sq])
			SetBit(&board.BlockerR315, R315[sq])
			board.Material[Black] += ValuePawn
		case 'n':
			SetBit(&board.Board[Black][Knight], sq)
			SetBit(&board.BlockerR90, R90[sq])
			SetBit(&board.BlockerR45, R45[sq])
			SetBit(&board.BlockerR315, R315[sq])
			board.Material[Black] += ValueKnight
		case 'b':
			SetBit(&board.Board[Black][Bishop], sq)
			SetBit(&board.BlockerR90, R90[sq])
			SetBit(&board.BlockerR45, R45[sq])
			SetBit(&board.BlockerR315, R315[sq])
			board.Material[Black] += ValueBishop
		case 'r':
			SetBit(&board.Board[Black][Rook], sq)
			SetBit(&board.BlockerR90, R90[sq])
			SetBit(&board.BlockerR45, R45[sq])
			SetBit(&board.BlockerR315, R315[sq])
			board.Material[Black] += ValueRook
		case 'q':
			SetBit(&board.Board[Black][Queen], sq)
			SetBit(&board.BlockerR90, R90[sq])
			SetBit(&board.BlockerR45, R45[sq])
			SetBit(&board.BlockerR315, R315[sq])
			board.Material[Black] += ValueQueen
		case 'k':
			SetBit(&board.Board[Black][King], sq)
			SetBit(&board.BlockerR90, R90[sq])
			SetBit(&board.BlockerR45, R45[sq])
			SetBit(&board.BlockerR315, R315[sq])
		case '/':
			r -= 8
			c = -1
		}

		if '0' <= char && char <= '9' {
			c += int(rune(char) - '0')
		} else {
			c++
		}

		if r == -8 && boardString[i+1] == ' ' {
			r = 0
		}

		if r < 0 || c > 8 {
			/* something went wrong */
			fmt.Printf("r: %d, c: %d\n", r, c)

			board = &ChessBoard{}
			return *board
		}

		if c == 8 && boardString[i+1] != '/' && boardString[i+1] != ' ' {
			/* something went wrong */
			board = &ChessBoard{}
			return *board
		}

		if char == ' ' {
			break
		}
	}

	board.PawnMaterial[White] = board.Material[White] - NumBits(board.Board[White][Pawn])*ValuePawn
	board.PawnMaterial[Black] = board.Material[Black] - NumBits(board.Board[Black][Pawn])*ValuePawn

	board.KingPos[White] = LeadBit(board.Board[White][King])
	board.KingPos[Black] = LeadBit(board.Board[Black][King])

	board.Friends[White] = board.Board[White][Pawn] | board.Board[White][Knight] | board.Board[White][Bishop] | board.Board[White][Rook] | board.Board[White][Queen] | board.Board[White][King]
	board.Friends[Black] = board.Board[Black][Pawn] | board.Board[Black][Knight] | board.Board[Black][Bishop] | board.Board[Black][Rook] | board.Board[Black][Queen] | board.Board[Black][King]
	board.Blocker = board.Friends[White] | board.Friends[Black]

	i++
	if len(boardString) <= i {
		/* something went wrong */
		board = &ChessBoard{}
		return *board
	} else if boardString[i] == 'w' {
		board.Side = White
		board.OurColor = White
	} else if boardString[i] == 'b' {
		board.Side = Black
		board.OurColor = Black
	} else {
		/* something went wrong */
		board = &ChessBoard{}
		return *board
	}
	i += 2
	for boardString[i] != ' ' {
		if boardString[i] == 'K' {
			board.CastleFlag |= WKingCastle
		} else if boardString[i] == 'Q' {
			board.CastleFlag |= WQueenCastle
		} else if boardString[i] == 'k' {
			board.CastleFlag |= BKingCastle
		} else if boardString[i] == 'q' {
			board.CastleFlag |= BQueenCastle
		} else if boardString[i] == '-' {
			i++
			break
		}
		i++
	}
	i++

	if boardString[i] != '-' {
		board.Ep = int16((boardString[i] - 'a') + (boardString[i+1]-'1')*8)
		i++
	} else {
		board.Ep = -1
	}
	return *board

	//Ignoring the rest of the EPD for now
}
