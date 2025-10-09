package graph

import (
    "context"
    "log"
    "sync"

    "go.uber.org/multierr"
)

type nodeOpt func(t *node)

func (n nodeOpt) Apply(t Any) {
    node, ok := t.(*node)
    if !ok {
        return
    }
    n(node)
}

type nodeData struct {
    data Any   // nil if error occurs
    err  error // if error occurs
}

// returns a node d with both fields is nil
func newNodeData() *nodeData {
    return &nodeData{}
}

// struct represent graph node
type node struct {
    rwLock    sync.RWMutex
    name      string     // name of this node, used for debug only
    inDegree  int64      // numbers of predecessors nodes of this one
    outDegree int64      // numbers of successors nodes of this one
    runner    Runner     // used to run this node
    d         *nodeData  // Data stores node result after execution
    g         *graphImpl // which graph this node belongs to
}

func withNodeName(name string) Opt {
    var f nodeOpt = func(t *node) {
        t.name = name
    }
    return f
}

func withGraph(g *graphImpl) Opt {
    var f nodeOpt = func(t *node) {
        t.g = g
    }
    return f
}

// return a node without it`s g uninitialized
func newNode(runner Runner) *node {
    return &node{runner: runner, d: newNodeData()}
}

func newNodeWithOpt(runner Runner, opts ...Opt) *node {
    n := newNode(runner)
    for _, opt := range opts {
        opt.Apply(n)
    }
    return n
}

// return a node it`s g initialized
func newBoundNode(runner Runner, g *graphImpl) *node {
    return newNodeWithOpt(runner, withGraph(g))
}

// return a named node
func newNamedNode(name string, runner Runner) *node {
    return newNodeWithOpt(runner, withNodeName(name))
}

func (n *node) fireExecution(ctx context.Context) {
    select {
    case <-ctx.Done():
        return
    case n.g.nodeReady <- n:
    }
}

func (n *node) decrementInDegree(ctx context.Context) {
    // check should directly return or not
    select {
    case <-ctx.Done():
        return
    default:
    }
    // decrementInDegree
    n.rwLock.Lock()
    n.inDegree = n.inDegree - 1
    n.rwLock.Unlock()
    // fire execution if in degree is equal to 0
    n.rwLock.RLock()
    defer n.rwLock.RUnlock()
    if n.inDegree == 0 {
        go n.fireExecution(ctx)
    }
}

func (n *node) decrementSuccessorInDegree(ctx context.Context) {
    for _, s := range n.g.successors[n] {
        go func(node *node) {
            node.decrementInDegree(ctx)
        }(s)
    }
}

func (n *node) finish(ctx context.Context) {
    select {
    case <-ctx.Done():
        return
    case n.g.nodeDone <- struct{}{}:
    }
}

func (n *node) Run(ctx context.Context, cancel context.CancelFunc) {
    // cas
    defer n.decrementSuccessorInDegree(ctx)
    defer n.finish(ctx)
    // from cache result
    if n.d.data != nil && n.d.err == nil {
        return
    }
    //
    errs := make([]error, 0)
    data := make([]Any, 0)
    for _, p := range n.g.predecessors[n] {
        data = append(data, p.d.data)
        errs = append(errs, p.d.err)
    }
    err := multierr.Combine(errs...)
    if err != nil {
        n.d.err = err
        cancel()
        return
    }
    n.d.data, n.d.err = n.runner.Run(ctx, data)
    if n.d.err != nil {
        cancel()
    }
}

// getNodeFromCtx fetch node from current context. For debug only.
func getNodeFromCtx(ctx context.Context) (n *node, ok bool) {
    n, ok = ctx.Value(&contextNodeInfoKey).(*node)
    return
}

// *graphImpl implements Graph
type graphImpl struct {
    done             chan struct{}              // a signal channel represent this graph is done or not
    nodeReady        chan *node                 // a signal channel used when a node is ready for execution
    nodeDone         chan struct{}              // a signal channel used when a node is done execution
    totalUnDoneNodes int                        // numbers of total undone nodes in this graph
    result           Result                     // cached result, lazily settled after graph had been successfully completed
    frozen           bool                       // whether this graph is frozen or not. the following fields should not be modified after graph is frozen
    successors       map[*node][]*node          // a node and its successors, e.g. n1 -> [n2,n3] represent edge from n1 to n2 and edge from n1 to n3
    predecessors     map[*node][]*node          // a node and its predecessors, e.g. n1 -> [n2, n3] represent edge from n2 to n1 and edge from n3 to n1
    totalNodes       map[DependAbleRunner]*node // a map used to store all the nodes in this graph
}

// Gather the result of nodes.
func (g *graphImpl) Gather(ctx context.Context) (Result, error) {
    g.mustFreeze()
    // return from cache
    if g.result != nil {
        result := g.result
        return result, nil
    }
    // check result from nodes
    errs := make([]error, 0)
    cache := make(map[DependAbleRunner]Any)
    for dr, n := range g.totalNodes {
        errs = append(errs, n.d.err)
        cache[dr] = n.d.data
    }
    err := multierr.Combine(errs...)
    if err != nil {
        return nil, err
    }
    // g.result is only settled after g has finished running without errors.
    g.result = cache
    return cache, nil
}

func (g *graphImpl) AddNodes(runner ...DependAbleRunner) {
    g.mustNotFreeze()
    for _, r := range runner {
        g.addNode(r)
    }
}

func (g *graphImpl) Freeze() {
    g.frozen = true
}

// Run is NOT GOROUTINE SAFE! Build multiple instance of graph if parallel execution is required.
// what is changed after a single iteration
// 1. node.inDegree is minus 1 until is equal to 0 for trigger execution of nodes, so we have to reset the InDegree of node in case of multiple call of Run
// 2. g.totalUnDoneNodes which will become the total number of nodes minus the number of nodes executed, therefore we need reset it too
func (g *graphImpl) Run(ctx context.Context) {
    g.mustFreeze()         // freeze from caller side
    g.resetUndoneCounter() // O(1)
    g.resetInDegree()      // O(V), actually quite simple to reset the counter
    ctxWithCancel, cancel := context.WithCancel(ctx)
    defer cancel()
    go g.selectInitReadyNodes(ctxWithCancel)
    go g.terminateWhenCountDown(ctxWithCancel)
    g.topologicalRun(ctxWithCancel, cancel)
}

// ////////////////////////////private methods////////////////////////
// check
var _ Graph = (*graphImpl)(nil)

func (g *graphImpl) mustFreeze() {
    if !g.frozen {
        panic("Graph should be frozen")
    }
}

func (g *graphImpl) mustNotFreeze() {
    if g.frozen {
        panic("Graph should not be frozen")
    }
}

func (g *graphImpl) resetUndoneCounter() {
    g.totalUnDoneNodes = len(g.totalNodes)
}

// reset in inDegree of node which has predecessors
func (g *graphImpl) resetInDegree() {
    for n, l := range g.predecessors {
        n.inDegree = int64(len(l))
    }
}

// try get or else create a node
// return true if already exists
func (g *graphImpl) getOrCreateNode(runner DependAbleRunner) (n *node, loaded bool) {
    n, loaded = g.totalNodes[runner]
    if loaded {
        return
    }
    n = newBoundNode(runner, g)
    g.totalNodes[runner] = n
    return
}

func (g *graphImpl) existsInSuccessor(pre, n *node) bool {
    successors, ok := g.successors[pre]
    if ok {
        for _, s := range successors {
            if s == n { // find an edge from pre to n
                return true
            }
        }
    }
    return false
}

// add edge if not exists
func (g *graphImpl) addEdgeIfNotExists(pre, n *node) (added bool) {
    // check successors
    if g.existsInSuccessor(pre, n) {
        return
    }
    g.successors[pre] = append(g.successors[pre], n)
    g.predecessors[n] = append(g.predecessors[n], pre)
    n.inDegree += 1
    pre.outDegree += 1
    added = true
    return
}

// add runner and its dependencies to this graph recursive
func (g *graphImpl) addNode(runner DependAbleRunner) {
    r := runner
    var addNodeRecursive func(runner DependAbleRunner, check bool)
    addNodeRecursive = func(runner DependAbleRunner, check bool) {
        if check && runner == r {
            panic("circle")
        }
        n, _ := g.getOrCreateNode(runner)
        dependencies := runner.GetDependency()
        if len(dependencies) > 0 {
            for _, pre := range dependencies {
                preNode, _ := g.getOrCreateNode(pre) // an edge from pre -> n
                added := g.addEdgeIfNotExists(preNode, n)
                if added { // cut recursive tree
                    addNodeRecursive(pre, true)
                }
            }
        }
    }
    addNodeRecursive(r, false)
}

func (g *graphImpl) terminateWhenCountDown(ctx context.Context) {
    defer func() {
        g.done <- struct{}{}
    }()
loop:
    for {
        select {
        case <-ctx.Done():
            return
        case <-g.nodeDone:
            g.totalUnDoneNodes -= 1
            if g.totalUnDoneNodes == 0 {
                break loop
            }
        }
    }
}

var contextNodeInfoKey int

func (g *graphImpl) topologicalRun(ctx context.Context, cancel context.CancelFunc) {
    defer cancel()
    for {
        select {
        case <-g.done:
            return
        case <-ctx.Done():
            log.Println("execution canceled due to cancellation of context")
            return
        case n := <-g.nodeReady:
            go func() {
                ctxWithNodeInfo := context.WithValue(ctx, &contextNodeInfoKey, n)
                n.Run(ctxWithNodeInfo, cancel)
            }()
        }
    }
}

func (g *graphImpl) selectInitReadyNodes(ctx context.Context) {
    for _, n := range g.totalNodes {
        // n has no predecessors
        if _, ok := g.predecessors[n]; !ok {
            go func(n *node) {
                select {
                case <-ctx.Done():
                    return
                case g.nodeReady <- n:
                }
            }(n)
        }
    }
}

// //////////////////// exported api //////////////////////

func NewGraph() Graph {
    g := &graphImpl{
        done:         make(chan struct{}),
        nodeReady:    make(chan *node),
        nodeDone:     make(chan struct{}),
        successors:   make(map[*node][]*node),
        predecessors: make(map[*node][]*node),
        totalNodes:   make(map[DependAbleRunner]*node),
    }
    return g
}
