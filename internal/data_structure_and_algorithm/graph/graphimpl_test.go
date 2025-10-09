package graph

import (
    "context"
    "errors"
    "fmt"
    "log"
    "math/rand"
    "testing"
    "time"
)

type d1 int

func (d *d1) Run(ctx context.Context, data []Any) (Any, error) {
    log.Println("node info")
    log.Println(ctx.Value(&contextNodeInfoKey))
    log.Println("begin run d1")
    log.Println(data)
    time.Sleep(time.Second)
    if fail := rand.Float64(); fail <= 0.5 { // fail at 50% percentage
        return "data1", errors.New("d1 error")
    }
    return "data1", nil
}

func (d *d1) GetDependency() []DependAbleRunner {
    return []DependAbleRunner{D4}
}

type d2 int

func (d *d2) Run(ctx context.Context, data []Any) (Any, error) {
    log.Println("node info")
    log.Println(ctx.Value(&contextNodeInfoKey))
    log.Println("begin run d2")
    log.Println(data)
    time.Sleep(time.Second)
    return "data2", nil
}

func (d *d2) GetDependency() []DependAbleRunner {
    return []DependAbleRunner{D1}
}

type d3 int

func (d *d3) Run(ctx context.Context, data []Any) (Any, error) {
    log.Println("node info")
    log.Println(ctx.Value(&contextNodeInfoKey))
    log.Println("begin run d3")
    log.Println(data)
    time.Sleep(time.Second)
    return "data3", nil
}

func (d *d3) GetDependency() []DependAbleRunner {
    return []DependAbleRunner{D2}
}

type d4 int

func (d *d4) Run(ctx context.Context, data []Any) (Any, error) {
    log.Println("node info")
    log.Println(ctx.Value(&contextNodeInfoKey))
    log.Println("begin run d4")
    log.Println(data)
    time.Sleep(time.Second)
    return "data4", nil
}

func (d *d4) GetDependency() []DependAbleRunner {
    return []DependAbleRunner{D1, D2, D3}
}

var d5 RunnerFn = func(ctx context.Context, data []Any) (Any, error) {
    log.Println("begin run d5")
    log.Println(data)
    return "data5", nil
}

var (
    D1 = NewRunnerBuilder().WithRunner(new(d1)).MustBuild()                         // nil -> d1
    D2 = NewRunnerBuilder().WithRunner(new(d2)).WithDependencies(D1).MustBuild()    // d1 -> d2
    D3 = NewRunnerBuilder().WithRunner(new(d3)).WithDependencies(D2).MustBuild()    // d2 -> d3
    D4 = BuildRunner(new(d4), D1, D2, D3)                                           // d1,d2,d3 -> d4
    D5 = NewRunnerBuilder().WithRunner(d5).WithDependencies(D2, D3, D1).MustBuild() // d1,d2,d3 -> d5
    D6 = NewRunnerBuilder().
        WithDependencies(D2, D3, D1).
            WithRunner(NewRunner(func(ctx context.Context, data []Any) (Any, error) {
                n, ok := getNodeFromCtx(ctx)
                if ok {
                    log.Println(n)
                }
                log.Println("begin run d6")
                log.Println(data)
                d2, d3, d1 := data[0], data[1], data[2]
                d2Result := d2.(string)
                d3Result := d3.(string)
                d1Result := d1.(string)
                s := fmt.Sprintf("data6: from %s-%s-%s", d2Result, d3Result, d1Result)
                log.Println("end run d6")
                log.Println(s)
                return s, nil
            })).MustBuild()
    // for testing circle
    // D1 = new(d1) // d4 -> d1
    // D2 = new(d2) // d1 -> d2
    // D3 = new(d3) // d2 -> d3
    // D4 = new(d4) // d1,d2,d3 -> d4
)

func TestBuildNode(t *testing.T) {
    n := newNamedNode("named", D5)
    t.Log(n.name)
}

func TestBuildNodeWithOption(t *testing.T) {
    n := newNamedNode("named", D5)
    opt := withNodeName("newNamed")
    opt.Apply(n)
    t.Log(n.name)
}

func TestGraphRun(t *testing.T) {
    ctx := context.Background()
    g := NewGraph()
    g.AddNodes(D6)
    g.Freeze()
    g.Run(ctx)
    g.Run(ctx)
    r, err := g.Gather(ctx)
    if err != nil {
        t.Log(err.Error())
        if r != nil {
            t.Fatal("map should be nil")
        }
    }
    t.Log(r)
}

func TestImplicitlyRun(t *testing.T) {
    r, err := WithImplicitlyGraph(D6).Run(context.Background())
    if err != nil {
        t.Fatal(err.Error())
    }
    t.Log(r)
}

func TestEmptyGraph(t *testing.T) {
    g := &graphImpl{
        done:         make(chan struct{}),
        nodeReady:    make(chan *node),
        nodeDone:     make(chan struct{}),
        successors:   make(map[*node][]*node),
        predecessors: make(map[*node][]*node),
        totalNodes:   make(map[DependAbleRunner]*node),
    }
    if g.result != nil {
        t.Fatal("g.result should be nil")
    }
}
