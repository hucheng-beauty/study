package registry

import (
    "context"
    "fmt"

    "github.com/jmoiron/sqlx"
    jsoniter "github.com/json-iterator/go"
)

var (
    _ Registry = SQLRegistry{}

    jsonCompatible = jsoniter.ConfigCompatibleWithStandardLibrary

    // when the sql result is too big, it's usually useless to print a large result to
    // see what's happening, we can just count on the result in %w at that time,
    // the large & useless error will give extra pressure to system's IO(not only
    // in log but also in potential response)
    // a useless example:
    // https://cloud.bytedance.net/argos/streamlog/info_overview/aggregate_statistics?
    // from=1685582477&pageType=location&
    // psm=lab.speech.saas&refresh=off&region=cn&to=1685625677
    // & it's been cut by log shower itself
    maxErrorPrintObjLength = 10000
)

type QueryBuilder interface {
    BuildInsert(tableName string, object any) (string, []any, error)
    BuildSelect(tableName string, object any, filter Filter) (string, []any, error)
    BuildUpdate(tableName string, object any, filter Filter) (string, []any, error)
    BuildDelete(tableName string, filter Filter) (string, []any, error)
}

type SQLRegistry struct {
    db           *sqlx.DB
    tableName    string
    queryBuilder QueryBuilder
}

func NewSQLRegistry(db *sqlx.DB, tableName string, queryBuilder QueryBuilder) SQLRegistry {
    return SQLRegistry{
        db:           db,
        tableName:    tableName,
        queryBuilder: queryBuilder,
    }
}

func (r SQLRegistry) Create(ctx context.Context, object any) (int64, error) {
    query, args, err := r.queryBuilder.BuildInsert(r.tableName, object)
    if err != nil {
        return 0, fmt.Errorf("build insert query: %w", err)
    }

    // log.V2.Info().With(ctx).
    //     Str("insert query:").Str(query).
    //     Str("object:").Obj(object).
    //     Str("args:").Obj(args).Emit()

    result, err := r.db.ExecContext(ctx, query, args...)
    if err != nil {
        return 0, r.formatErr("exec insert", query, args, object, nil, err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("get last inserted id: %w", err)
    }
    return id, nil
}

func (r SQLRegistry) Update(ctx context.Context, object any, filter Filter) error {
    query, args, err := r.queryBuilder.BuildUpdate(r.tableName, object, filter)
    if err != nil {
        return fmt.Errorf("build partial update query: %w", err)
    }

    result, err := r.db.ExecContext(ctx, query, args...)
    if err != nil {
        return r.formatErr("exec update", query, args, object, filter, err)
    }

    rowsAffected, rowsErr := result.RowsAffected()
    if rowsErr != nil {
        rowsAffected = -420
    }

    fmt.Println(rowsAffected)
    // log.V2.Info().With(ctx).
    //     Str("update query:").Str(query).
    //     Str("rows affected:").Int64(rowsAffected).
    //     Str("object:").Obj(object).
    //     Str("args:").Obj(args).Emit()
    return nil
}

func (r SQLRegistry) Select(ctx context.Context, object any, filter Filter) error {
    query, args, err := r.queryBuilder.BuildSelect(r.tableName, object, filter)
    if err != nil {
        return fmt.Errorf("build select: %w", err)
    }

    err = r.db.GetContext(ctx, object, query, args...)
    if err != nil {
        return r.formatErr("get", query, args, object, filter, err)
    }

    // NOTE: get query is too frequent, disable log for now.
    // log.V2.Info().With(ctx).
    //     Str("get query:").Str(query).
    //     Str("table_name:").Str(r.tableName).
    //     Str("filter:").Obj(filter).
    //     Str("object:").Obj(object).
    //     Str("args:").Obj(args).Emit()

    return nil
}

func (r SQLRegistry) Describe(ctx context.Context, objects any, filter Filter) error {
    query, args, err := r.queryBuilder.BuildSelect(r.tableName, objects, filter)
    if err != nil {
        return fmt.Errorf("build select: %w", err)
    }

    // log.V2.Info().With(ctx).
    //     Str("select query:").Str(query).
    //     Str("table_name:").Str(r.tableName).
    //     Str("filter:").Obj(filter).Emit()

    err = r.db.SelectContext(ctx, objects, query, args...)
    if err != nil {
        return r.formatErr("select", query, args, objects, filter, err)
    }
    return nil
}

func (r SQLRegistry) Delete(ctx context.Context, filter Filter) error {
    query, args, err := r.queryBuilder.BuildDelete(r.tableName, filter)
    if err != nil {
        return fmt.Errorf("build delete: %w", err)
    }

    _, err = r.db.ExecContext(ctx, query, args...)
    if err != nil {
        return r.formatErr("exec delete", query, args, nil, filter, err)
    }

    // log.V2.Info().With(ctx).
    //     Str("delete query:").Str(query).
    //     Str("args:").Obj(args).
    //     Str("table_name:").Str(r.tableName).
    //     Str("filter:").Obj(filter).Emit()

    return nil
}

// formatErr fills error with data that's needed to debug errors in SQL calls.
// special logic: when the serialization of obj is too large,
// we'll limit the output size to reduce I/O lose
func (r SQLRegistry) formatErr(callType, query string, args []any, object any, filter Filter, err error) error {
    argStr, _ := jsonCompatible.MarshalToString(args)
    objStr, _ := jsonCompatible.MarshalToString(object)
    filterStr, _ := jsonCompatible.MarshalToString(filter)
    objStrLen := len(objStr)

    if objStrLen <= maxErrorPrintObjLength {
        return fmt.Errorf("[query=%s, args=%s, object=%s, filter=%s] %s: %w",
            query, argStr, objStr, filterStr, callType, err)
    }

    obj := objStr[:maxErrorPrintObjLength]
    return fmt.Errorf("[query=%s, args=%s, object=%s ...(about more %d length not shown), filter=%s] %s: %w",
        query, argStr, obj, objStrLen-maxErrorPrintObjLength, filterStr, callType, err)
}
