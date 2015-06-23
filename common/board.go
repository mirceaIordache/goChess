package common

import "github.com/op/go-logging"

const MaxPlyDepth = 20

type BitBoard uint64

var ChessLogger *logging.Logger

type ChessBoard struct {
	Board [2][7]BitBoard /*piece position by side (0-white, 1-black)
	  and then piece (1-pawn..6-king) */
	Friends      [2]BitBoard /* Friendly pieces */
	Blocker      BitBoard
	BlockerR90   BitBoard  /* Rotated 90deg */
	BlockerR45   BitBoard  /* Rotated 45deg */
	BlockerR315  BitBoard  /* Rotated 315deg */
	Ep           int16     /* En Passant Square */
	CastleFlag   int16     /* Flag for castle privileges */
	Side         uint16    /* Side to move, 0-white, 1-black */
	Material     [2]uint16 /* Total material by side */
	PawnMaterial [2]uint16 /* Total pawn material by side */
	KingPos      [2]uint16 /* King Positions 0-A1..63-h8 */
	Castled      [2]bool   /* Side if castled */
	OurColor     uint16    /* Our color */
	GameCount    int
}

type Move struct {
	Move  int /* The move itself */
	Score int /* Scored value of the move */
}
