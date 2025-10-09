```go

import (
    "context"
    "common/graph"
)

func TestGraph(t *testing.T) {
    g := graph.NewGraph()
    g.AddNodes()
    g.Freeze()
    g.Run(context.TODO())
}
```
