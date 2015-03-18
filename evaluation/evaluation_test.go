package evaluation

import (
	"github.com/mirceaIordache/goChess/common"

	"testing"
)

func TestEvaluation(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"8/8/3K4/8/5n2/8/7k/8 w - -", 0},                      // Draw
		{"8/8/3k4/8/3K4/4P3/8/8 w - -", 100},                   // KPK
		{"8/8/3k4/8/3K4/4P3/8/8 b - -", -100},                  // KPK
		{"8/8/3Kp3/8/5n2/1r6/7k/8 w - - ", -1000},              // Lone King
		{"8/8/3Kp3/8/5n2/1r6/7k/8 b - - ", 1000},               // Lone King
		{"8/5b2/3K4/8/5n2/8/7k/8 b - -", 700},                  // KBNK
		{"8/5b2/3K4/8/5n2/8/7k/8 w - -", -700},                 // KBNK
		{"2r3k1/pppR1pp1/4p3/4P1P1/5P2/1P4K1/P1P5/8 w - -", 0}, // BK6 -- Yes, should be 0, it's a BK test.
	}

	for _, c := range cases {
		board := common.FromEPD(c.in)
		got := 0
		if EvaluateDraw(board) == false {
			got = Evaluate(32767, -32767, board)
		}

		if got != c.want {
			t.Errorf("IN: %q, WANT: %d, OUT: %d,", c.in, c.want, got)
		}
	}
}
