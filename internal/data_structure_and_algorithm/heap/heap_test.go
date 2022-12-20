package heap

import (
	"testing"
)

func TestMainer(t *testing.T) {
	h := NewHeap(10)
	h.Insert(5)
	h.Insert(6)

	t.Logf("%d\n", h.Top())
	h.RemoveTop()
	t.Logf("%d\n", h.Top())
}
