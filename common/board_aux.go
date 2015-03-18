package common

/* Set i-th bit in the bitboard */
func SetBit(b *BitBoard, i int) {
	*b |= BitPosArray[i]
}

/* Clear i-th bit in the bitboard */
func ClearBit(b *BitBoard, i int) {
	*b = *b & NotBitPosArray[i]
}

/* Get number of 1 bits in the bitboard */
func NumBits(b BitBoard) uint16 {
	return BitCountArray[b>>48] + BitCountArray[(b>>32)&0xffff] + BitCountArray[(b>>16)&0xffff] + BitCountArray[b&0xffff]
}

/* Get the leading 1 bit in the bitboard */
/* TODO Make it efficient */
func LeadBit(b BitBoard) uint16 {
	var i uint16
	for i = 64; i >= 1; i-- {
		if b>>i != 0 {
			return i
		}
	}

	return 0
}

func TrailBit(b BitBoard) uint16 {
	return LeadBit(b & (^b + 1))
}

/* Generate a CBoard
CBoard - a 64 element array with all pieces, colour-neutral */
func GenerateCBoard(board ChessBoard) [64]uint16 {
	var b BitBoard
	var sq, piece uint16
	var cboard [64]uint16

	for piece = Pawn; piece <= King; piece++ {
		b = board.Board[White][piece] | board.Board[Black][piece]
		for b != 0 {
			sq = LeadBit(b)
			ClearBit(&b, int(sq))
			cboard[sq] = piece
		}
	}

	return cboard
}

func GenerateMoveBoard(board ChessBoard) [64]uint16 {
	var moveBoard [64]uint16

	return moveBoard
}

/* Calculate the material (dis)advantage for the side */
func CalculateMaterial(board ChessBoard) int {
	return int(board.Material[board.Side]) - int(board.Material[1^board.Side])
}

func CalculatePieces(board ChessBoard, side uint16) BitBoard {
	return board.Board[side][Knight] | board.Board[side][Bishop] | board.Board[side][Rook] | board.Board[side][Queen]
}

/* Calculate the game phase based on total material */
/* For more information on why there are more phases than usual (opening, mid-game, end-game)
   see https://chessprogramming.wikispaces.com/Evaluation+Discontinuity */

func GetGamePhase(board ChessBoard) uint16 {
	return 8 - (board.Material[White]+board.Material[Black])/1150
}
