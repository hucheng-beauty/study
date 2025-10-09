package graph

import "context"

type Any = interface{}

// Runner is a basic extension point for every logic
type Runner interface {
    Run(ctx context.Context, data []Any) (Any, error)
}

// RunnerFn implements Runner
// a help function for building a Runner
type RunnerFn func(ctx context.Context, data []Any) (Any, error)

func (r RunnerFn) Run(ctx context.Context, data []Any) (Any, error) {
    return r(ctx, data)
}

func NewRunner(f RunnerFn) Runner {
    return f
}

// DependAbleRunner is runner with dependencies
type DependAbleRunner interface {
    Runner
    GetDependency() []DependAbleRunner
}

// RunnerBuilder is the builder for DependAbleRunner
type RunnerBuilder interface {
    WithRunner(r Runner) RunnerBuilder
    // WithDependencies returns RunnerBuilder with dependant of a list of runners, order matters.
    WithDependencies(dependencies ...DependAbleRunner) RunnerBuilder
    MustBuild() DependAbleRunner
}

// NewRunnerBuilder build runner with builder pattern
func NewRunnerBuilder() RunnerBuilder {
    return &defaultDependAbleRunner{}
}

// BuildRunner build runner directly
func BuildRunner(runner Runner, dependencies ...DependAbleRunner) DependAbleRunner {
    return &defaultDependAbleRunner{r: runner, dependencies: dependencies}
}

// Opt of graph or node
type Opt interface {
    // Apply this option to t
    Apply(t Any)
}

type Result map[DependAbleRunner]Any

func (r Result) Get(dr DependAbleRunner) (v Any, ok bool) {
    v, ok = r[dr]
    return
}

// Graph interface exported
type Graph interface {
    // AddNodes add a list of runner to this graph
    // be cautioned if d is not in runners but in the dependencies of some runner in runners, d will be implicitly added to this graph
    AddNodes(runners ...DependAbleRunner)
    // Freeze  this graph after all nodes and edges is added to this graph. Freeze should be called before Run.
    Freeze()
    // Run fire execution in topological sort order
    Run(ctx context.Context)
    // Gather collects all nodes result after run. all result should be returned as nil if graph has not been fired yet.
    Gather(ctx context.Context) (Result, error)
}

// private
type defaultDependAbleRunner struct {
    r            Runner
    dependencies []DependAbleRunner
}

func (d *defaultDependAbleRunner) WithRunner(r Runner) RunnerBuilder {
    d.r = r
    return d
}

func (d *defaultDependAbleRunner) WithDependencies(dependencies ...DependAbleRunner) RunnerBuilder {
    d.dependencies = dependencies
    return d
}

func (d *defaultDependAbleRunner) MustBuild() DependAbleRunner {
    if d.r == nil {
        panic("nil runner")
    }
    return d
}

func (d *defaultDependAbleRunner) Run(ctx context.Context, data []Any) (Any, error) {
    return d.r.Run(ctx, data)
}

func (d *defaultDependAbleRunner) GetDependency() []DependAbleRunner {
    return d.dependencies
}

// helper struct for building implicitly graph
type withImplicitlyGraph struct {
    dr DependAbleRunner
}

func (ig *withImplicitlyGraph) Run(ctx context.Context) (Any, error) {
    g := NewGraph()
    g.AddNodes(ig.dr)
    g.Freeze()
    g.Run(ctx)
    r, err := g.Gather(ctx)
    if err != nil {
        return nil, err
    }
    return r[ig.dr], nil
}

func WithImplicitlyGraph(dr DependAbleRunner) *withImplicitlyGraph {
    return &withImplicitlyGraph{dr}
}
