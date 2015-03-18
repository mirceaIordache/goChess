package common

/* colours */
const (
	White = iota
	Black = iota
)

/* pieces */
const (
	Empty  = iota
	Pawn   = iota
	Knight = iota
	Bishop = iota
	Rook   = iota
	Queen  = iota
	King   = iota
	BPawn  = iota
)

var PawnType = [2]int{Pawn, BPawn}

const Notation = " PNBRQK"
const LNotation = " pnbrqk"

var Algebraic = [...]string{"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
	"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
	"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
	"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
	"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
	"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
	"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
	"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8"}

/* board indices */
const (
	A1 = iota
	B1 = iota
	C1 = iota
	D1 = iota
	E1 = iota
	F1 = iota
	G1 = iota
	H1 = iota

	A2 = iota
	B2 = iota
	C2 = iota
	D2 = iota
	E2 = iota
	F2 = iota
	G2 = iota
	H2 = iota

	A3 = iota
	B3 = iota
	C3 = iota
	D3 = iota
	E3 = iota
	F3 = iota
	G3 = iota
	H3 = iota

	A4 = iota
	B4 = iota
	C4 = iota
	D4 = iota
	E4 = iota
	F4 = iota
	G4 = iota
	H4 = iota

	A5 = iota
	B5 = iota
	C5 = iota
	D5 = iota
	E5 = iota
	F5 = iota
	G5 = iota
	H5 = iota

	A6 = iota
	B6 = iota
	C6 = iota
	D6 = iota
	E6 = iota
	F6 = iota
	G6 = iota
	H6 = iota

	A7 = iota
	B7 = iota
	C7 = iota
	D7 = iota
	E7 = iota
	F7 = iota
	G7 = iota
	H7 = iota

	A8 = iota
	B8 = iota
	C8 = iota
	D8 = iota
	E8 = iota
	F8 = iota
	G8 = iota
	H8 = iota
)

const (
	AFile = iota
	BFile = iota
	CFile = iota
	DFile = iota
	EFile = iota
	FFile = iota
	GFile = iota
	HFile = iota
)

/* material values */
const ValuePawn = 100
const ValueKnight = 350
const ValueBishop = 350
const ValueRook = 550
const ValueQueen = 1100
const ValueKing = 2000

var Values = [...]uint16{0, 100, 350, 350, 550, 1100, 2000, 0}

/* Node Types */
const (
	PV  = iota
	ALL = iota
	CUT = iota
)

/* castling flags */
const WKingCastle = 0x1
const WQueenCastle = 0x2
const BKingCastle = 0x4
const BQueenCastle = 0x8
const WCastle = WKingCastle | WQueenCastle
const BCastle = BKingCastle | BQueenCastle

/* Mode Detection */
const KnightPromotion = 0x00002000
const BishopPromotion = 0x00003000
const RookPromotion = 0x00004000
const QueenPromotion = 0x00005000
const PiecePromotion = 0x00007000
const PawnCapture = 0x00008000
const KnightCapture = 0x00010000
const BishopCapture = 0x00018000
const RookCapture = 0x00020000
const QueenCapture = 0x00028000
const PieceCapture = 0x00038000
const NullMove = 0x00100000
const CastlingMove = 0x00200000
const EnPassantMove = 0x00300000
const MoveMask = (CastlingMove | EnPassantMove | PiecePromotion | 0x0FFF)

/* Mate detection */
const Mate = 32767

/* some special bitboards */
const NullBitBoard = BitBoard(0x0000000000000000)
const WhiteSquares = BitBoard(0x55AA55AA55AA55AA)
const BlackSquares = BitBoard(0xAA55AA55AA55AA55)
const CentreSquares = BitBoard(0x0000001818000000)
const ComputerHalf = BitBoard(0xFFFFFFFF00000000)
const OpponentHalf = BitBoard(0x00000000FFFFFFFF)

/* bitboard arrays */
var BitPosArray = [...]BitBoard{
	0x1, 0x2, 0x4, 0x8,
	0x10, 0x20, 0x40, 0x80,
	0x100, 0x200, 0x400, 0x800,
	0x1000, 0x2000, 0x4000, 0x8000,
	0x10000, 0x20000, 0x40000, 0x80000,
	0x100000, 0x200000, 0x400000, 0x800000,
	0x1000000, 0x2000000, 0x4000000, 0x8000000,
	0x10000000, 0x20000000, 0x40000000, 0x80000000,
	0x100000000, 0x200000000, 0x400000000, 0x800000000,
	0x1000000000, 0x2000000000, 0x4000000000, 0x8000000000,
	0x10000000000, 0x20000000000, 0x40000000000, 0x80000000000,
	0x100000000000, 0x200000000000, 0x400000000000, 0x800000000000,
	0x1000000000000, 0x2000000000000, 0x4000000000000, 0x8000000000000,
	0x10000000000000, 0x20000000000000, 0x40000000000000, 0x80000000000000,
	0x100000000000000, 0x200000000000000, 0x400000000000000, 0x800000000000000,
	0x1000000000000000, 0x2000000000000000, 0x4000000000000000, 0x8000000000000000}

var NotBitPosArray = [...]BitBoard{
	0xFFFFFFFFFFFFFFFE, 0xFFFFFFFFFFFFFFFD, 0xFFFFFFFFFFFFFFFB, 0xFFFFFFFFFFFFFFF7,
	0xFFFFFFFFFFFFFFEF, 0xFFFFFFFFFFFFFFDF, 0xFFFFFFFFFFFFFFBF, 0xFFFFFFFFFFFFFF7F,
	0xFFFFFFFFFFFFFEFF, 0xFFFFFFFFFFFFFDFF, 0xFFFFFFFFFFFFFBFF, 0xFFFFFFFFFFFFF7FF,
	0xFFFFFFFFFFFFEFFF, 0xFFFFFFFFFFFFDFFF, 0xFFFFFFFFFFFFBFFF, 0xFFFFFFFFFFFF7FFF,
	0xFFFFFFFFFFFEFFFF, 0xFFFFFFFFFFFDFFFF, 0xFFFFFFFFFFFBFFFF, 0xFFFFFFFFFFF7FFFF,
	0xFFFFFFFFFFEFFFFF, 0xFFFFFFFFFFDFFFFF, 0xFFFFFFFFFFBFFFFF, 0xFFFFFFFFFF7FFFFF,
	0xFFFFFFFFFEFFFFFF, 0xFFFFFFFFFDFFFFFF, 0xFFFFFFFFFBFFFFFF, 0xFFFFFFFFF7FFFFFF,
	0xFFFFFFFFEFFFFFFF, 0xFFFFFFFFDFFFFFFF, 0xFFFFFFFFBFFFFFFF, 0xFFFFFFFF7FFFFFFF,
	0xFFFFFFFEFFFFFFFF, 0xFFFFFFFDFFFFFFFF, 0xFFFFFFFBFFFFFFFF, 0xFFFFFFF7FFFFFFFF,
	0xFFFFFFEFFFFFFFFF, 0xFFFFFFDFFFFFFFFF, 0xFFFFFFBFFFFFFFFF, 0xFFFFFF7FFFFFFFFF,
	0xFFFFFEFFFFFFFFFF, 0xFFFFFDFFFFFFFFFF, 0xFFFFFBFFFFFFFFFF, 0xFFFFF7FFFFFFFFFF,
	0xFFFFEFFFFFFFFFFF, 0xFFFFDFFFFFFFFFFF, 0xFFFFBFFFFFFFFFFF, 0xFFFF7FFFFFFFFFFF,
	0xFFFEFFFFFFFFFFFF, 0xFFFDFFFFFFFFFFFF, 0xFFFBFFFFFFFFFFFF, 0xFFF7FFFFFFFFFFFF,
	0xFFEFFFFFFFFFFFFF, 0xFFDFFFFFFFFFFFFF, 0xFFBFFFFFFFFFFFFF, 0xFF7FFFFFFFFFFFFF,
	0xFEFFFFFFFFFFFFFF, 0xFDFFFFFFFFFFFFFF, 0xFBFFFFFFFFFFFFFF, 0xF7FFFFFFFFFFFFFF,
	0xEFFFFFFFFFFFFFFF, 0xDFFFFFFFFFFFFFFF, 0xBFFFFFFFFFFFFFFF, 0x7FFFFFFFFFFFFFFF}

/* rotations */

var Boxes = [2]BitBoard{0x00003C3C3C3C0000, 0x007E7E7E7E7E7E00}
var Stonewall = [2]BitBoard{0x81400000000, 0x14080000}
var Rings = [4]BitBoard{0x0000001818000000, 0x00003C24243C0000, 0x007E424242427E00, 0xFF818181818181FF}
var Rank6 = [2]uint16{5, 2}
var Rank7 = [2]uint16{6, 1}
var Rank8 = [2]uint16{7, 0}
var MaskKR_TrappedWhite = [3]BitBoard{0x1000000000000, 0x101000000000000, 0x301000000000000}
var MaskQR_TrappedWhite = [3]BitBoard{0x80000000000000, 0x8080000000000000, 0xc080000000000000}
var MaskKR_TrappedBlack = [3]BitBoard{0x100, 0x101, 0x103}
var MaskQR_TrappedBlack = [3]BitBoard{0x8000, 0x8080, 0x80c0}
var BoardHalf = [2]BitBoard{0xffffffff00000000, 0x00000000ffffffff}
var BoardSide = [2]BitBoard{0x0f0f0f0f0f0f0f0f, 0xf0f0f0f0f0f0f0f0}
var Slider = [8]int{0, 0, 0, 1, 1, 1, 0, 0}
