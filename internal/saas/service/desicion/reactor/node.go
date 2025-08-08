package reactor

import (
    "context"
    "fmt"

    "study/internal/saas/model"
    "study/internal/saas/toolkit/tree"

    "github.com/samber/lo"
)

type (
    predicator func(ctx context.Context, oc *orderContext) (bool, error)
    executor   func(ctx context.Context, oc *orderContext) error
)

type value interface {
    fmt.Stringer
    Predicate(ctx context.Context, oc *orderContext) (bool, error)
}

type node tree.Node[value]

type orderContext struct {
    order *model.Order
    // instance  *trading.Instance
    // collector orderCollector.Collector

    // optional means that the struct in here is not available
    // until manually injected
    optional
}

type optional struct {
    // instanceRegistry volcanoInstance.VolcanoPlantInstanceRegistry

    // prepaidPlantAttributes is used only by enablePrepaidService as a shortcut to avoid parse attr again.
    // prepaidPlantAttributes *prepaidPlantAttributes

    // resourcePackWithoutPlantAttributes is used only by enableResourcePackStandalone as a shortcut to avoid
    // parse attr again.
    // resourcePackWithoutPlantAttributes *resourcePackWithoutPlantAttributes

    // formalizeResourcePackAttributes is used only by enableFormalizeResourcePack as a shortcut to
    // avoid parse attr again.
    // formalizeResourcePackAttributes *formalizeResourcePackAttributes

    service *model.Plant

    // subOrder                        *trading.SubOrder
    subOrder interface{}

    // nonSystemResource includes all non-system
    // resource with the instance number from order.
    nonSystemResource []model.BundleResourcePackage
}

type actionValue struct {
    name      string
    predicate predicator
    executor  executor
}

func newAction(name string, predicate predicator, executor executor) node {
    return tree.NewSimpleNode(
        value(&actionValue{
            name:      name,
            predicate: predicate,
            executor:  executor,
        }),
    )
}

func (r *actionValue) String() string { return r.name }

func (r *actionValue) Predicate(ctx context.Context, oc *orderContext) (bool, error) {
    return r.predicate(ctx, oc)
}

type strategyValue struct {
    name       string
    predicator predicator
}

func newStrategy(name string, predicate predicator, nodes ...node) node {
    // Note: useful for debugging, unreachable when running
    if len(nodes) <= 0 {
        panic("strategy should always have at least one child node")
    }

    simpleNode := tree.NewSimpleNode(
        value(&strategyValue{
            name:       name,
            predicator: predicate,
        }),
    )
    simpleNode.SetChildren(lo.Map(nodes,
        func(n node, _ int) tree.Node[value] { return tree.Node[value](n) })...,
    )

    return simpleNode
}

func (r *strategyValue) String() string { return r.name }

func (r *strategyValue) Predicate(ctx context.Context, oc *orderContext) (bool, error) {
    return r.predicator(ctx, oc)
}

func passAll(_ context.Context, _ *orderContext) (bool, error) { return true, nil }

func not(predicate predicator) predicator {
    return func(ctx context.Context, oc *orderContext) (bool, error) {
        ok, err := predicate(ctx, oc)
        if err != nil {
            return false, err
        }

        return !ok, nil
    }
}
