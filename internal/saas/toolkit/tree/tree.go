package tree

type Node[T any] interface {
    GetValue() T
    GetChildren() []Node[T]
    SetChildren(...Node[T])
    GetParent() Node[T]
    SetParent(Node[T])
}

func NewSimpleNode[T any](value T) Node[T] {
    return &simple[T]{value: value, children: nil}
}

type simple[T any] struct {
    value    T
    children []Node[T]
    parent   Node[T]
}

func (s *simple[T]) GetValue() T { return s.value }

func (s *simple[T]) GetChildren() []Node[T] { return s.children }

func (s *simple[T]) SetChildren(nodes ...Node[T]) {
    s.setChildren(true, nodes...)
}

func (s *simple[T]) setChildren(firstCall bool, nodes ...Node[T]) {
    for _, node := range nodes {
        if node != s {
            s.children = append(s.children, node)

            n, ok := node.(*simple[T])
            if firstCall && ok {
                n.setParent(false, s)
            }
        }
    }
}

func (s *simple[T]) GetParent() Node[T] { return s.parent }

func (s *simple[T]) SetParent(node Node[T]) {
    s.setParent(true, node)
}

func (s *simple[T]) setParent(firstCall bool, node Node[T]) {
    s.parent = node

    n, ok := node.(*simple[T])
    if firstCall && ok {
        n.setChildren(false, s)
    }
}
