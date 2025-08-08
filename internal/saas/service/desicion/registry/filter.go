package registry

import (
    "fmt"
    "strings"
)

const (
    andSeparator = " AND "
    orSeparator  = " OR "

    TrueLogicalConst  = "TRUE"
    FalseLogicalConst = "FALSE"
)

var _ OrderingFilter = sliceFilter{}

type sliceFilter struct {
    filters []Filter

    sqlSeparator   string
    statefulFilter func(state, filterMatches bool) bool
}

func (sf sliceFilter) SQL() (string, []any) {
    var (
        bracedQueries []string
        args          []any
    )

    for _, filter := range sf.filters {
        filterQuery, filterArgs := filter.SQL()
        bracedQueries = append(bracedQueries, fmt.Sprintf("(%s)", filterQuery))
        args = append(args, filterArgs...)
    }

    query := strings.Join(bracedQueries, sf.sqlSeparator)
    return query, args
}

func (sf sliceFilter) Matches(o any) (bool, error) {
    state, err := sf.filters[0].Matches(o)
    if err != nil {
        return false, fmt.Errorf("[filter_index=0] check match: %w", err)
    }

    for i, filter := range sf.filters[1:] {
        match, err := filter.Matches(o)
        if err != nil {
            return false, fmt.Errorf("[filter_index=%v] check match: %w", i+1, err)
        }

        state = sf.statefulFilter(state, match)
    }

    return state, nil
}

func (sf sliceFilter) SQLOrderBy() string {
    var orderComponents []string

    for _, filter := range sf.filters {
        of, ok := filter.(OrderingFilter)
        if !ok {
            continue
        }

        orderBy := of.SQLOrderBy()
        if orderBy != "" {
            orderComponents = append(orderComponents, orderBy)
        }
    }

    if len(orderComponents) > 0 {
        return strings.Join(orderComponents, OrderSeparator)
    }

    return ""
}

func (sf sliceFilter) Less(object1, object2 any) (bool, error) {
    for _, filter := range sf.filters {
        of, ok := filter.(OrderingFilter)
        if !ok {
            continue
        }

        // if object1 < object2 return true
        lessLeft, errL := of.Less(object1, object2)
        if errL != nil {
            return false, errL
        }
        if lessLeft {
            return true, nil
        }

        // if object1 > object2 return false
        lessRight, errG := of.Less(object2, object1)
        if errG != nil {
            return false, errG
        }
        if lessRight {
            return false, nil
        }

        // !(object1 > object2 || object1 < object2) => object1 == object2
        // this implies we need to check using next filter
    }

    return false, nil
}

// Or combines one or more filters using OR logical operator into a new Filter.
func Or(filter Filter, filters ...Filter) Filter {
    return sliceFilter{
        filters: append([]Filter{filter}, filters...),

        sqlSeparator:   orSeparator,
        statefulFilter: func(state, filterMatches bool) bool { return state || filterMatches },
    }
}

// And combines one or more filters using AND logical operator into a new Filter.
func And(filter Filter, filters ...Filter) Filter {
    return sliceFilter{
        filters: append([]Filter{filter}, filters...),

        sqlSeparator:   andSeparator,
        statefulFilter: func(state, filterMatches bool) bool { return state && filterMatches },
    }
}

type allFilter struct{}

// All returns a Filter that matches all entities.
func All() Filter { return allFilter{} }

// SQL for allFilter a TRUE logical constant.
func (allFilter) SQL() (string, []any) { return TrueLogicalConst, nil }

// Matches for allFilter always returns true.
func (allFilter) Matches(any) (bool, error) { return true, nil }

type noneFilter struct{}

// None returns a Filter that matches no entity.
func None() Filter { return noneFilter{} }

// SQL for noneFilter is a FALSE logical constant.
func (noneFilter) SQL() (string, []any) { return FalseLogicalConst, nil }

// Matches for noneFilter always returns false.
func (noneFilter) Matches(object any) (bool, error) { return false, nil }
