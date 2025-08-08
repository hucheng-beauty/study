package reactor

import (
    "context"
    "errors"
    "fmt"
    "strings"

    jsoniter "github.com/json-iterator/go"
)

var (
    jsonCompatible          = jsoniter.ConfigCompatibleWithStandardLibrary
    errRegisterPathNotFound = errors.New("found no register branch")
)

type dfsInfo struct {
    ctx          context.Context
    orderContext *orderContext
    path         []string
    executor
    err error
}

type register struct{ root node }

func (r *register) MatchAndExecute(ctx context.Context, oc *orderContext) ([]string, error) {
    info := &dfsInfo{ctx: ctx, orderContext: oc}

    if isStopped := r.dfs(r.root, info); !isStopped {
        return nil, errRegisterPathNotFound
    }
    if info.err != nil {
        path := joinPath(info.path)
        errP := fmt.Errorf("[register] path: %s, predicate error: %w",
            path, info.err)
        return nil, errP
    }

    err := info.executor(ctx, info.orderContext)
    if err != nil {
        path := joinPath(info.path)
        errE := fmt.Errorf("[register] path: %s, executor error: %w",
            path, err)
        return nil, errE
    }
    return info.path, nil
}

func (r *register) dfs(node node, info *dfsInfo) (isStopped bool) {
    if node == nil {
        return false
    }

    nodeValue := node.GetValue()
    if nodeValue == nil {
        return false
    }

    ok, err := nodeValue.Predicate(info.ctx, info.orderContext)
    info.path = append(info.path, nodeValue.String())
    defer func() {
        if !isStopped {
            info.path = info.path[:len(info.path)-1]
        }
    }()
    if err != nil {
        info.err = err
        return true
    }

    if !ok {
        return false
    }

    action, isAction := nodeValue.(*actionValue)
    if isAction {
        info.executor = action.executor
        return true
    }

    for _, ch := range node.GetChildren() {
        isStopped = r.dfs(ch, info)
        if isStopped {
            return true
        }
    }

    return false
}

func joinPath(path []string) string { return strings.Join(path, "/") }

func and(predicators ...predicator) predicator {
    return func(ctx context.Context, oc *orderContext) (bool, error) {
        for _, predicate := range predicators {
            ok, err := predicate(ctx, oc)
            if err != nil {
                return false, err
            }

            if !ok {
                return false, nil
            }
        }

        return true, nil
    }
}
