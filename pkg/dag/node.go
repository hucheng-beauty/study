package dag

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Any = interface{}

// Runner is basic extension point for every logic
type Runner interface {
	Run(ctx context.Context, data []Any) (Any, error)
}

type node struct {
	mutex     sync.RWMutex // mutex ensures concurrency security.
	name      string       // name is the name of runner.
	in        []Any        // in is Runner's input parameter of the data.
	out       Any          // out is thr result of runner.
	err       error        // err is thr error of runner.
	inDegree  int          // inDegree is the number of adjacencyList's value of the graph.
	outDegree int          // outDegree is the number of adjacencyList's key of the graph.
	r         Runner
	g         *graph
}

func (n *node) Run(ctx context.Context, in []Any) (Any, error) {
	result, err := n.r.Run(ctx, in)
	n.out = result
	n.err = err
	return result, err
}

type NodeOption func(*node)

func withName(name string) NodeOption {
	return func(n *node) {
		n.name = name
	}
}

func withIn(in []Any) NodeOption {
	return func(n *node) {
		n.in = in
	}
}

func withOut(out Any) NodeOption {
	return func(n *node) {
		n.out = out
	}
}

func withErr(err error) NodeOption {
	return func(n *node) {
		n.err = err
	}
}

func withRunner(r Runner) NodeOption {
	return func(n *node) {
		n.r = r
	}
}

func withGraph(g *graph) NodeOption {
	return func(n *node) {
		n.g = g
	}
}

// NewNodeWithOption create a node using option.
func NewNodeWithOption(opts ...NodeOption) *node {
	n := new(node)
	for _, opt := range opts {
		opt(n)
	}
	return n
}

// DependAbleRunner is runner with dependencies
type DependAbleRunner interface {
	Runner
	GetDependency() []DependAbleRunner
}

type dependAbleNode struct {
	n *node
	d []DependAbleRunner
}

func (d *dependAbleNode) Run(ctx context.Context, in []Any) (Any, error) {
	return d.n.r.Run(ctx, in)
}

func (d *dependAbleNode) GetDependency() []DependAbleRunner {
	return d.d
}

type DependAbleNodeOption func(*dependAbleNode)

func withNode(n *node) DependAbleNodeOption {
	return func(dn *dependAbleNode) {
		dn.n = n
	}
}

func withDependAbleRunner(dr ...DependAbleRunner) DependAbleNodeOption {
	return func(dn *dependAbleNode) {
		dn.d = append(dn.d, dr...)
	}
}

// NewDependAbleNode create a node with dependencies using option.
func NewDependAbleNode(opts ...DependAbleNodeOption) *dependAbleNode {
	n := new(dependAbleNode)
	for _, opt := range opts {
		opt(n)
	}
	return n
}

// TODO: The data is aim to test.
type d1 int

func (d *d1) Run(ctx context.Context, data []Any) (Any, error) {
	log.Println("begin run d1")
	log.Printf("d1 indata: %+v\n", data)
	time.Sleep(time.Second)
	if fail := rand.Float64(); fail <= 0.5 { // fail at 50% percentage
		return "data1", errors.New("d1 error")
	}
	return "data1", nil
}

func (d *d1) GetDependency() []DependAbleRunner {
	return []DependAbleRunner{}
}

type d2 int

func (d *d2) Run(ctx context.Context, data []Any) (Any, error) {
	log.Println("begin run d2")
	log.Printf("d2 indata: %+v\n", data)
	time.Sleep(time.Second)
	return "data2", nil
}

func (d *d2) GetDependency() []DependAbleRunner {
	return []DependAbleRunner{D1}
}

type d3 int

func (d *d3) Run(ctx context.Context, data []Any) (Any, error) {
	log.Println("begin run d3")
	log.Printf("d3 indata:%+v\n", data)
	time.Sleep(time.Second)
	return "data3", nil
}

func (d *d3) GetDependency() []DependAbleRunner {
	return []DependAbleRunner{D2}
}

type d4 int

func (d *d4) Run(ctx context.Context, data []Any) (Any, error) {
	log.Println("begin run d4")
	log.Printf("d4 indata:%+v\n", data)
	time.Sleep(time.Second)
	return "data4", nil
}

func (d *d4) GetDependency() []DependAbleRunner {
	return []DependAbleRunner{D2, D3}
}

var (
	D1 = new(d1)
	D2 = new(d2)
	D3 = new(d3)
	D4 = new(d4)
)
