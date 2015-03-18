package evaluation

import . "github.com/mirceaIordache/goChess/common"

const SCORE_HUNG = -20
const SCORE_DEV_KINGMOVED = -20
const SCORE_DEV_ROOKMOVED = -20
const SCORE_DEV_NOTCASTLED = -8
const SCORE_TRADE_PIECE = 4
const SCORE_TRADE_PAWN = 8
const SCORE_PFACTOR = 550
const SCORE_DRAW = 0

const SCORE_PAWN_BLOCKED = -48
const SCORE_PAWN_BASEATTACK = -12
const SCORE_PAWN_EIGHT = -10
const SCORE_PAWN_STONEWALL = -10
const SCORE_PAWN_LOCKED = -10
const SCORE_PAWN_EARLYWINGMOVE = -6
const SCORE_PAWN_EARLYCENTERPREPEAT = -6
const SCORE_PAWN_CONNECTED = 50
const SCORE_PAWN_NEARKING = 40
const SCORE_PAWN_ATACKWEAK = 2
const SCORE_PAWN_FIANCHETTO = -13

const SCORE_KNIGHT_PINNED = -30
const SCORE_KNIGHT_ONRIM = -13
const SCORE_KNIGHT_OUTPOST = 10

const SCORE_BISHOP_TRAPPED = -250
const SCORE_BISHOP_PINNED = -30
const SCORE_BISHOP_FIANCHETTO = 8
const SCORE_BISHOP_OUTPOST = 8
const SCORE_BISHOP_DOUBLE = 18

const SCORE_ROOK_PINNED = -40
const SCORE_ROOK_TRAPPED = -10
const SCORE_ROOK_HALFFILE = 5
const SCORE_ROOK_OPENFILE = 6
const SCORE_ROOK_7RANK = 30
const SCORE_ROOK_LIBERATED = 40

const SCORE_QUEEN_PINNED = -90
const SCORE_QUEEN_EARLYMOVE = -40
const SCORE_QUEEN_ABSENT = -25
const SCORE_QUEEN_NEARKING = 12

const SCORE_KING_HOPEN = -300
const SCORE_KING_DEFENDERDEFICIT = -50
const SCORE_KING_BACK_RANK_WEAK = -40
const SCORE_KING_GOPEN = -30
const SCORE_KING_RUPTURE = -20
const SCORE_KING_OPENFILE = -10
const SCORE_KING_ENEMYOPENFILE = -6

var pawnScores = [2][64]int{
	{0, 0, 0, 0, 0, 0, 0, 0,
		5, 5, 5, -10, -10, 5, 5, 5,
		-2, -2, -2, 6, 6, -2, -2, -2,
		0, 0, 0, 25, 25, 0, 0, 0,
		2, 2, 12, 16, 16, 12, 2, 2,
		4, 8, 12, 16, 16, 12, 4, 4,
		4, 8, 12, 16, 16, 12, 4, 4,
		0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0,
		4, 8, 12, 16, 16, 12, 4, 4,
		4, 8, 12, 16, 16, 12, 4, 4,
		2, 2, 12, 16, 16, 12, 2, 2,
		0, 0, 0, 25, 25, 0, 0, 0,
		-2, -2, -2, 6, 6, -2, -2, -2,
		5, 5, 5, -10, -10, 5, 5, 5,
		0, 0, 0, 0, 0, 0, 0, 0},
}

var passedPawnScore = [2][8]int{{0, 48, 48, 120, 144, 192, 240, 0}, {0, 240, 192, 144, 120, 48, 48, 0}}
var IsolaniNormal = [8]int{-8, -10, -12, -14, -14, -12, -10, -8}
var IsolaniWeaker = [8]int{-22, -24, -26, -28, -28, -26, -24, -22}
var pawnCover = [9]int{-60, -30, 0, 5, 30, 30, 30, 30, 30}
var factor = [9]int{7, 8, 8, 7, 6, 5, 4, 2, 0}

var d2e2 = [2]BitBoard{0x0018000000000000, 0x0000000000001800}
var brank7 = [2]BitBoard{0x000000000000FF00, 0x00FF000000000000}
var brank67 = [2]BitBoard{0x0000000000FFFF00, 0x00FFFF0000000000}
var brank8 = [2]BitBoard{0x00000000000000FF, 0xFF00000000000000}
var Outpost = [2][64]int{
	{0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 1, 0, 0,
		0, 0, 0, 1, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 0, 0,
		0, 1, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0},
}
var KingSq = [64]int{
	24, 24, 24, 16, 16, 0, 32, 32,
	24, 20, 16, 12, 12, 16, 20, 24,
	16, 12, 8, 4, 4, 8, 12, 16,
	12, 8, 4, 0, 0, 4, 8, 12,
	12, 8, 4, 0, 0, 4, 8, 12,
	16, 12, 8, 4, 4, 8, 12, 16,
	24, 20, 16, 12, 12, 16, 20, 24,
	24, 24, 24, 16, 16, 0, 32, 32,
}
var EndingKing = [64]int{
	0, 6, 12, 18, 18, 12, 6, 0,
	6, 12, 18, 24, 24, 18, 12, 6,
	12, 18, 24, 32, 32, 24, 18, 12,
	18, 24, 32, 48, 48, 32, 24, 18,
	18, 24, 32, 48, 48, 32, 24, 18,
	12, 18, 24, 32, 32, 24, 18, 12,
	6, 12, 18, 24, 24, 18, 12, 6,
	0, 6, 12, 18, 18, 12, 6, 0,
}

var KBNK = [64]int{
	0, 10, 20, 30, 40, 50, 60, 70,
	10, 20, 30, 40, 50, 60, 70, 60,
	20, 30, 40, 50, 60, 70, 60, 50,
	30, 40, 50, 60, 70, 60, 50, 40,
	40, 50, 60, 70, 60, 50, 40, 30,
	50, 60, 70, 60, 50, 40, 30, 20,
	60, 70, 60, 50, 40, 30, 20, 10,
	70, 60, 50, 40, 30, 20, 10, 0}

var nn = [2]BitBoard{0x4200000000000000, 0x0000000000000042}
var bb = [2]BitBoard{0x2400000000000000, 0x0000000000000024}
