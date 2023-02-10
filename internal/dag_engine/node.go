package dag_engine

import (
	"context"
	"sync"
)

type Node struct {
	ctx     context.Context
	mutex   sync.RWMutex
	Name    string
	inData  []Any
	outData Any
	err     error
	r       Runner
}

func (n *Node) Run(ctx context.Context, in []Any) (Any, error) {
	r, err := n.r.Run(ctx, in)
	n.outData = r
	n.err = err
	return r, err
}

func NewNode(opts ...NodeOption) *Node {
	n := &Node{
		Name:    "",
		ctx:     context.Background(),
		outData: nil,
	}

	for _, opt := range opts {
		opt(n)
	}
	return n
}

type NodeOption func(*Node)

func WithCtx(ctx context.Context) NodeOption {
	return func(node *Node) {
		node.ctx = ctx
	}
}

func WithName(name string) NodeOption {
	return func(node *Node) {
		node.Name = name
	}
}

func WithInData(inData []Any) NodeOption {
	return func(node *Node) {
		node.outData = inData
	}
}

func WithOutData(outData Any) NodeOption {
	return func(node *Node) {
		node.outData = outData
	}
}

func WithError(err error) NodeOption {
	return func(node *Node) {
		node.err = err
	}
}

func WithRunner(r Runner) NodeOption {
	return func(node *Node) {
		node.r = r
	}
}

type DependAbleNode struct {
	Node *Node
	d    []DependAbleRunner
}

func (d *DependAbleNode) Run(ctx context.Context, inData []Any) (Any, error) {
	return d.Node.r.Run(ctx, inData)
}

func (d *DependAbleNode) GetDependency() []DependAbleRunner {
	return d.d
}

func NewDependAbleNode(opts ...DependAbleNodeOption) *DependAbleNode {
	n := &DependAbleNode{
		Node: nil,
		d:    nil,
	}

	for _, opt := range opts {
		opt(n)
	}
	return n
}

type DependAbleNodeOption func(*DependAbleNode)

func WithNode(n *Node) DependAbleNodeOption {
	return func(node *DependAbleNode) {
		node.Node = n
	}
}

func WithDependAbleRunner(d []DependAbleRunner) DependAbleNodeOption {
	return func(node *DependAbleNode) {
		node.d = append(node.d, d...)
	}
}
