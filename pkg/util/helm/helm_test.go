package helm

import "testing"

func TestPageOffset(t *testing.T) {
	cases := []struct {
		pageNum int
		want    int
	}{
		{pageNum: 0, want: 0},
		{pageNum: -3, want: 0},
		{pageNum: 1, want: 0},
		{pageNum: 2, want: 1},
		{pageNum: 10, want: 9},
	}
	for _, c := range cases {
		if got := pageOffset(c.pageNum); got != c.want {
			t.Errorf("pageOffset(%d) = %d, want %d", c.pageNum, got, c.want)
		}
	}
}
