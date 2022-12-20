package dag

import (
	"context"
	"testing"
)

func TestMainer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g := NewGraph()
	g.AddNodes(D4)
	g.Run(ctx, cancel)
	g.Run(ctx, cancel)
}
