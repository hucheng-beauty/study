package dag_engine

import "context"

type Any = interface{}

// Runner basic extension point for every logic
type Runner interface {
	Run(ctx context.Context, data []Any) (Any, error)
}

// DependAbleRunner runner with dependencies
type DependAbleRunner interface {
	Runner
	GetDependency() []DependAbleRunner
}

type Container interface {
	Set(key DependAbleRunner, value *DependAbleNode)
	Get(key DependAbleRunner) *DependAbleNode
}

type AdjLister interface {
	GetInstance() map[DependAbleRunner][]DependAbleRunner
	Set(key DependAbleRunner, value []DependAbleRunner)
	Get(key DependAbleRunner) []DependAbleRunner
}

type Counter interface {
	Add(number int)
	Sub(number int)
	Get() int
}

type Device interface {
	Set(key DependAbleRunner, value int)
	Get(key DependAbleRunner) int
}

type Visitor interface {
	Set(key string, value bool)
	Get(key string) bool
}
