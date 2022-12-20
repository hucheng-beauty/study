package dag_engine

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
)

/*
	1. run
	2. node
	3. yaml
	4. 字段名称 ==》 逻辑
*/

type dag struct {
	Count int
	dnm   Container
	ra    AdjLister
	wq    chan DependAbleRunner
	end   chan struct{}
	retry int32
}

func (d *dag) addMap(dRunner DependAbleRunner, dNode *DependAbleNode) {
	d.dnm.Set(dRunner, dNode)
}

func (d *dag) addRunner(dRunner DependAbleRunner) {
	d.ra.Set(dRunner, dRunner.GetDependency())
}

func (d *dag) Run(ctx context.Context, cancel context.CancelFunc) {
	count := NewCount(0)
	var retry int32

	log.Println("begin generate adjacency list")
	adjList := getAdjList(d.ra)

	log.Println("begin initialize the in-degree")
	inDegree := getInDegree(adjList)

	log.Println("begin check isRing")
	if isRing(d.ra, d.Count) {
		return
	}

	log.Println("begin run dag engine")
	for key, _ := range adjList.GetInstance() {
		if inDegree.Get(key) == 0 {
			d.wq <- key
		}
	}

	go func() {
		for {
			if len(d.wq) > 0 {
				popNode := <-d.wq

				go func(popNode DependAbleRunner) {
					out := make([]Any, 0)
					if len(d.ra.Get(popNode)) > 0 {
						for _, r := range d.dnm.Get(popNode).d {
							out = append(out, d.dnm.Get(r).Node.outData)
						}
					}

					_, err := d.dnm.Get(popNode).Node.Run(nil, out)
					if err != nil {
						log.Printf("[dag_engine][Run] error: %s\n", err.Error())
						cancel()
					}
					count.Add(1)

					select {
					case <-ctx.Done():
						return
					default:
						for _, node := range adjList.Get(popNode) {
							fmt.Printf("%+#v, %d\n", node, inDegree.Get(node))
							newInDegree := inDegree.Get(node) - 1
							inDegree.Set(node, newInDegree)
							fmt.Printf("%+#v, %d\n", node, inDegree.Get(node))
							if inDegree.Get(node) == 0 {
								d.wq <- node
							}
						}
					}
				}(popNode)
			}
		}
	}()

	go func() {
		for {
			if count.Get() == d.Count {
				fmt.Println("dag engine do success")
				count = NewCount()
				inDegree = getInDegree(adjList)
				atomic.AddInt32(&retry, 1)
				d.end <- struct{}{}
			} else if count.Get() == d.Count && retry == d.retry {
				fmt.Println("dag engine do success")
				d.end <- struct{}{}
			}
		}
	}()

	select {
	case <-d.end:
		return
	case <-ctx.Done():
		return
	}
}

func NewDAG(count int) *dag {
	return &dag{
		Count: count,
		dnm:   NewDependAbleNodeMap(),
		ra:    NewAdjList(),
		wq:    make(chan DependAbleRunner, count),
		end:   make(chan struct{}),
		retry: 2,
	}
}

func getAdjList(rAdjList AdjLister) AdjLister {
	adjList := NewAdjList()
	for ka, va := range rAdjList.GetInstance() {
		for _, vva := range va {
			sli := adjList.Get(vva)
			sli = append(sli, ka)
			adjList.Set(vva, sli)
		}
	}
	return adjList
}

func getInDegree(adjList AdjLister) Device {
	i := NewInDegree()
	for _, v := range adjList.GetInstance() {
		for _, node := range v {
			i.Set(node, i.Get(node)+1)
		}
	}
	return i
}

func isRing(ra AdjLister, count int) bool {
	c := NewCount()
	adjList := getAdjList(ra)
	i := getInDegree(adjList)

	queue := make([]DependAbleRunner, 0)
	for key, _ := range adjList.GetInstance() {
		if i.Get(key) == 0 {
			queue = append(queue, key)
		}
	}

	for len(queue) > 0 {
		popNode := queue[0]
		queue = queue[1:]
		c.Add(1)

		for _, node := range adjList.Get(popNode) {
			i.Set(node, i.Get(node)-1)
			if i.Get(node) == 0 {
				queue = append(queue, node)
			}
		}
	}

	if c.Get() != count && len(queue) >= 0 {
		log.Println("有环")
		return true
	}
	return false
}

func runnerProcess(in []DependAbleRunner) (out []DependAbleRunner) {
	if len(in) <= 0 {
		return in
	}

	for _, runner := range in {
		if runner.GetDependency() != nil {
			out = append(out, runnerProcess(runner.GetDependency())...)
		}
		out = append(out, runner)
	}
	return Deduplicate(out)
}

func BuildGraphFromRunners(runners ...DependAbleRunner) *dag {
	ar := make([]DependAbleRunner, 0)
	for _, runner := range runners {
		ar = append(ar, runner)
	}
	rs := runnerProcess(ar)

	dag := NewDAG(len(rs))

	for _, runner := range rs {
		dn := NewDependAbleNode(WithNode(NewNode(WithRunner(runner))), WithDependAbleRunner(runner.GetDependency()))
		dag.addRunner(runner)
		dag.addMap(runner, dn)
	}

	return dag
}
