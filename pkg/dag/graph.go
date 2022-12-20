package dag

import (
	"context"
)

type graph struct {
	done           chan struct{}                           // done is end flag
	nq             chan *node                              // nq just is a chan of store node.
	count          int                                     // count is the number of node in graph
	frozen         bool                                    // frozen can judge whether the graph stops.
	adjList        map[DependAbleRunner][]DependAbleRunner // adjList, AdjacencyList: [n1 -> n2], [n1 -> n3].
	inverseAdjList map[DependAbleRunner][]DependAbleRunner // inverseAdjList, InverseAdjacencyList: [n2 -> n1], [n3 -> n1].
	dnm            map[DependAbleRunner]*dependAbleNode    // dnm is the container of store the relationship of DependAbleRunner and node.
}

func (g *graph) AddNodes(drs ...DependAbleRunner) {
	if g.frozen {
		panic("[dag][AddNode] Graph should not be frozen.")
	}

	for _, dr := range drs {
		g.addNode(dr)
		g.addEdge(dr)
	}
}

func (g *graph) addNode(dr DependAbleRunner) {
	dr, exist := g.dnm[dr]
	if exist {
		panic("generating ring")
	}

	n := withNode(NewNodeWithOption(withRunner(dr)))
	dn := NewDependAbleNode(n, withDependAbleRunner(dr.GetDependency()...))
	g.dnm[dr] = dn
}

func (g *graph) addEdge(dr DependAbleRunner) {
	g.adjList[dr] = dr.GetDependency()
}

func (g *graph) Run(ctx context.Context, cancel context.CancelFunc) {
	select {
	case <-ctx.Done():
		return
	case n := <-g.nq:
		if n.inDegree == 0 {
			n.out, n.err = n.Run(ctx, n.in)
		} else {

		}
	}
}

func NewGraph() *graph {
	return &graph{
		done:           make(chan struct{}),
		nq:             make(chan *node),
		count:          -1,
		frozen:         false,
		adjList:        make(map[DependAbleRunner][]DependAbleRunner),
		inverseAdjList: make(map[DependAbleRunner][]DependAbleRunner),
		dnm:            make(map[DependAbleRunner]*dependAbleNode),
	}
}
