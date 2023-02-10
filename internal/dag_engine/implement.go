package dag_engine

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

type DependAbleNodeMap struct {
	mutex sync.RWMutex
	m     map[DependAbleRunner]*DependAbleNode // map[DependAbleRunner]DependAbleNode
}

func (dnm *DependAbleNodeMap) Set(key DependAbleRunner, value *DependAbleNode) {
	if key == nil {
		return
	}

	dnm.mutex.Lock()
	defer dnm.mutex.Unlock()

	dnm.m[key] = value
}

func (nm *DependAbleNodeMap) Get(key DependAbleRunner) *DependAbleNode {
	if key == nil {
		return &DependAbleNode{}
	}

	nm.mutex.RLock()
	defer nm.mutex.RUnlock()

	return nm.m[key]
}

func (nm *DependAbleNodeMap) Len() int {
	nm.mutex.RLock()
	defer nm.mutex.RUnlock()

	return len(nm.m)
}

func NewDependAbleNodeMap() *DependAbleNodeMap {
	return &DependAbleNodeMap{
		mutex: sync.RWMutex{},
		m:     make(map[DependAbleRunner]*DependAbleNode),
	}
}

type AdjList struct {
	sync.RWMutex
	m map[DependAbleRunner][]DependAbleRunner // map[node_name]dependants
}

func (al *AdjList) GetInstance() map[DependAbleRunner][]DependAbleRunner {
	al.RWMutex.RLock()
	defer al.RWMutex.RUnlock()

	return al.m
}

func (al *AdjList) Set(key DependAbleRunner, value []DependAbleRunner) {
	al.RWMutex.Lock()
	defer al.RWMutex.Unlock()

	al.m[key] = value
}

func (al *AdjList) Get(key DependAbleRunner) []DependAbleRunner {
	al.RWMutex.RLock()
	defer al.RWMutex.RUnlock()

	return al.m[key]
}

func NewAdjList() *AdjList {
	return &AdjList{
		RWMutex: sync.RWMutex{},
		m:       make(map[DependAbleRunner][]DependAbleRunner),
	}
}

type InDegree struct {
	sync.RWMutex
	m map[DependAbleRunner]int // map[DependAbleRunner]inDegree
}

func (i *InDegree) Set(key DependAbleRunner, value int) {
	i.RWMutex.Lock()
	defer i.RWMutex.Unlock()

	i.m[key] = value
}

func (i *InDegree) Get(key DependAbleRunner) int {
	i.RWMutex.RLock()
	defer i.RWMutex.RUnlock()

	return i.m[key]
}

func NewInDegree() *InDegree {
	return &InDegree{
		RWMutex: sync.RWMutex{},
		m:       make(map[DependAbleRunner]int),
	}
}

type Count struct {
	sync.RWMutex
	n int
}

func (c *Count) Add(number int) {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	c.n += number
}

func (c *Count) Sub(number int) {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	c.n -= number
}

func (c *Count) Get() int {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()

	return c.n
}

func NewCount(count ...int) *Count {
	c := &Count{
		RWMutex: sync.RWMutex{},
		n:       0,
	}

	if count != nil {
		c.n = count[0]
	}
	return c
}

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
