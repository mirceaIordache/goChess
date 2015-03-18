package common

var PassedPawnMask = [2][64]BitBoard{{0xc0c0c0c0c0c0c0,
	0xe0e0e0e0e0e0e0, 0x70707070707070, 0x38383838383838, 0x1c1c1c1c1c1c1c, 0xe0e0e0e0e0e0e, 0x7070707070707, 0x3030303030303, 0xc0c0c0c0c0c0, 0xe0e0e0e0e0e0, 0x707070707070,
	0x383838383838, 0x1c1c1c1c1c1c, 0xe0e0e0e0e0e, 0x70707070707, 0x30303030303, 0xc0c0c0c0c0, 0xe0e0e0e0e0, 0x7070707070, 0x3838383838, 0x1c1c1c1c1c,
	0xe0e0e0e0e, 0x707070707, 0x303030303, 0xc0c0c0c0, 0xe0e0e0e0, 0x70707070, 0x38383838, 0x1c1c1c1c, 0xe0e0e0e, 0x7070707,
	0x3030303, 0xc0c0c0, 0xe0e0e0, 0x707070, 0x383838, 0x1c1c1c, 0xe0e0e, 0x70707, 0x30303, 0xc0c0,
	0xe0e0, 0x7070, 0x3838, 0x1c1c, 0xe0e, 0x707, 0x303, 0xc0, 0xe0, 0x70,
	0x38, 0x1c, 0xe, 0x7, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0}, {0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc000000000000000, 0xe000000000000000, 0x7000000000000000,
	0x3800000000000000, 0x1c00000000000000, 0xe00000000000000, 0x700000000000000, 0x300000000000000, 0xc0c0000000000000, 0xe0e0000000000000, 0x7070000000000000, 0x3838000000000000, 0x1c1c000000000000,
	0xe0e000000000000, 0x707000000000000, 0x303000000000000, 0xc0c0c00000000000, 0xe0e0e00000000000, 0x7070700000000000, 0x3838380000000000, 0x1c1c1c0000000000, 0xe0e0e0000000000, 0x707070000000000,
	0x303030000000000, 0xc0c0c0c000000000, 0xe0e0e0e000000000, 0x7070707000000000, 0x3838383800000000, 0x1c1c1c1c00000000, 0xe0e0e0e00000000, 0x707070700000000, 0x303030300000000, 0xc0c0c0c0c0000000,
	0xe0e0e0e0e0000000, 0x7070707070000000, 0x3838383838000000, 0x1c1c1c1c1c000000, 0xe0e0e0e0e000000, 0x707070707000000, 0x303030303000000, 0xc0c0c0c0c0c00000, 0xe0e0e0e0e0e00000, 0x7070707070700000,
	0x3838383838380000, 0x1c1c1c1c1c1c0000, 0xe0e0e0e0e0e0000, 0x707070707070000, 0x303030303030000, 0xc0c0c0c0c0c0c000, 0xe0e0e0e0e0e0e000, 0x7070707070707000, 0x3838383838383800, 0x1c1c1c1c1c1c1c00,
	0xe0e0e0e0e0e0e00, 0x707070707070700, 0x303030303030300}}

var SquarePawnMask = [2][64]BitBoard{{0xffffffffffffffff,
	0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xfcfcfcfcfcfc, 0xfefefefefefe, 0xffffffffffff,
	0xffffffffffff, 0xffffffffffff, 0xffffffffffff, 0x7f7f7f7f7f7f, 0x3f3f3f3f3f3f, 0xfcfcfcfcfcfc, 0xfefefefefefe, 0xffffffffffff, 0xffffffffffff, 0xffffffffffff,
	0xffffffffffff, 0x7f7f7f7f7f7f, 0x3f3f3f3f3f3f, 0xf8f8f8f8f8, 0xfcfcfcfcfc, 0xfefefefefe, 0xffffffffff, 0xffffffffff, 0x7f7f7f7f7f, 0x3f3f3f3f3f,
	0x1f1f1f1f1f, 0xf0f0f0f0, 0xf8f8f8f8, 0xfcfcfcfc, 0xfefefefe, 0x7f7f7f7f, 0x3f3f3f3f, 0x1f1f1f1f, 0xf0f0f0f, 0xe0e0e0,
	0xf0f0f0, 0xf8f8f8, 0x7c7c7c, 0x3e3e3e, 0x1f1f1f, 0xf0f0f, 0x70707, 0xc0c0, 0xe0e0, 0x7070,
	0x3838, 0x1c1c, 0xe0e, 0x707, 0x303, 0x80, 0x40, 0x20, 0x10, 0x8,
	0x4, 0x2, 0x1}, {0x8000000000000000,
	0x4000000000000000, 0x2000000000000000, 0x1000000000000000, 0x800000000000000, 0x400000000000000, 0x200000000000000, 0x100000000000000, 0xc0c0000000000000, 0xe0e0000000000000, 0x7070000000000000,
	0x3838000000000000, 0x1c1c000000000000, 0xe0e000000000000, 0x707000000000000, 0x303000000000000, 0xe0e0e00000000000, 0xf0f0f00000000000, 0xf8f8f80000000000, 0x7c7c7c0000000000, 0x3e3e3e0000000000,
	0x1f1f1f0000000000, 0xf0f0f0000000000, 0x707070000000000, 0xf0f0f0f000000000, 0xf8f8f8f800000000, 0xfcfcfcfc00000000, 0xfefefefe00000000, 0x7f7f7f7f00000000, 0x3f3f3f3f00000000, 0x1f1f1f1f00000000,
	0xf0f0f0f00000000, 0xf8f8f8f8f8000000, 0xfcfcfcfcfc000000, 0xfefefefefe000000, 0xffffffffff000000, 0xffffffffff000000, 0x7f7f7f7f7f000000, 0x3f3f3f3f3f000000, 0x1f1f1f1f1f000000, 0xfcfcfcfcfcfc0000,
	0xfefefefefefe0000, 0xffffffffffff0000, 0xffffffffffff0000, 0xffffffffffff0000, 0xffffffffffff0000, 0x7f7f7f7f7f7f0000, 0x3f3f3f3f3f3f0000, 0xfcfcfcfcfcfc0000, 0xfefefefefefe0000, 0xffffffffffff0000,
	0xffffffffffff0000, 0xffffffffffff0000, 0xffffffffffff0000, 0x7f7f7f7f7f7f0000, 0x3f3f3f3f3f3f0000, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff,
	0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff}}