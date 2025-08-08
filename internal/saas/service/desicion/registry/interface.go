package registry

import "context"

type Registry interface {
    Loader
    Modifier
}

type Loader interface {
    Select(ctx context.Context, object any, filter Filter) error
    Describe(ctx context.Context, object any, filter Filter) error
}

type Modifier interface {
    Create(ctx context.Context, object any, ) (int64, error)
    Update(ctx context.Context, object any, filter Filter) error
    Delete(ctx context.Context, filter Filter) error
}

// Filter return entity that can filter data returned from registries.
type Filter interface {
    SQL() (query string, args []any)
    Matches(object any) (bool, error)
}
