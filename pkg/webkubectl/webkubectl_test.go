package webkubectl

import "testing"

func TestName(t *testing.T) {
	w := NewWebkubectl()
	if err := w.Run(); err != nil {
		t.Error(err)
	}
}
