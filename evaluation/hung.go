package evaluation

import . "github.com/mirceaIordache/goChess/common"

func HungPieces(board ChessBoard, side uint16, attackBoard [2][7]BitBoard) int {
	xside := 1 ^ side
	hunged := 0

	/* Knight */
	n := attackBoard[xside][Pawn] & board.Board[side][Knight]
	n |= (attackBoard[xside][0] & board.Board[side][Knight] & ^attackBoard[side][0])

	/* Bishop */
	b := attackBoard[xside][Pawn] & board.Board[side][Bishop]
	b |= (attackBoard[xside][0] & board.Board[side][Bishop] & ^attackBoard[side][0])

	/* Rook */

	r := attackBoard[xside][Pawn] | attackBoard[xside][Knight] | attackBoard[xside][Bishop]
	r &= board.Board[side][Rook]
	r |= (attackBoard[xside][0] & board.Board[side][Rook] & ^attackBoard[side][0])

	/* Queen */
	q := attackBoard[xside][Pawn] | attackBoard[xside][Knight] | attackBoard[xside][Bishop] | attackBoard[xside][Rook]
	q &= board.Board[side][Queen]
	q |= (attackBoard[xside][0] & board.Board[side][Queen] & ^attackBoard[side][0])

	c := n | b | r | q

	if c != 0 {
		hunged = int(NumBits(c))
	}

	/* King */
	if attackBoard[xside][0]&board.Board[side][King] != 0 {
		hunged++
	}

	return hunged
}
