package main

import "testing"

func TestMakeBoard(t *testing.T) {
	b := MakeBoard(5, 3)
	if b.Rows() != 5 || b.Cols() != 3 {
		t.Fail()
	}
}

func TestBoardAccess(t *testing.T) {
	b := MakeBoard(3, 4)
	b.Set(0, 0, true)
	b.Set(2, 3, true)
	b.Set(1, 1, true)

	test := []struct {
		r, c int
		want bool
	}{
		{0, 0, true},
		{0, 1, false},
		{0, 2, false},
		{0, 3, false},
		{1, 0, false},
		{1, 1, true},
		{1, 2, false},
		{1, 3, false},
		{2, 0, false},
		{2, 1, false},
		{2, 2, false},
		{2, 3, true},

		{0, 4, true},
		{-1, 0, false},
		{0, -1, false},
		{3, 0, true},
	}

	for _, c := range test {
		if have := b.Get(c.r, c.c); have != c.want {
			t.Errorf("%#v: got: %v", c, have)
		}
	}
}
