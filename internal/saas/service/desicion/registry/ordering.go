package registry

import "fmt"

const (
    // OrderDirectionIgnore is a special case of OrderDirection in
    // which no ordering by given field should be done.
    OrderDirectionIgnore OrderDirection = ""
    // OrderDirectionAscending orders by field in ascending order.
    OrderDirectionAscending OrderDirection = "ASC"
    // OrderDirectionDescending orders by field in descending order.
    OrderDirectionDescending OrderDirection = "DESC"

    OrderSeparator = ", "
)

type OrderDirection string

func (od OrderDirection) SQL() string {
    switch od {
    case OrderDirectionAscending:
        return string(OrderDirectionAscending)
    case OrderDirectionDescending:
        return string(OrderDirectionDescending)

    default:
        panic(fmt.Sprintf("cannot generate SQL for order direction %s", od))
    }
}

func (od OrderDirection) Less(object1, object2 any) (bool, error) {
    result, err := compareBuiltin(object1, object2)
    if err != nil {
        return false, fmt.Errorf("[object1=%+#v, object2=%+#v] compareBuiltin: %w",
            object1, object2, err)
    }

    switch od {
    case OrderDirectionAscending:
        return result == compareResultLessThan, nil
    case OrderDirectionDescending:
        return result == compareResultGreaterThan, nil

    default:
        return false, fmt.Errorf("comparison invalid for order direction %s", od)
    }
}

type Order interface {
    // SQLOrderBy returns ORDER BY clause to be used in SQL SELECT query for ordering results.
    SQLOrderBy() string
    // Less return true if object1 should be returned before object2 in response.
    Less(object1, object2 any) (bool, error)
}

// OrderingFilter extends Filter interface functionality by adding an option to order query results.
type OrderingFilter interface {
    Filter
    Order
}

type orderingFilter struct {
    filter Filter
    order  Order
}

// ApplyOrder combines Order with Filter to produce OrderingFilter.
func ApplyOrder(filter Filter, order Order) OrderingFilter {
    return orderingFilter{filter: filter, order: order}
}

func (of orderingFilter) SQL() (string, []any) { return of.filter.SQL() }

func (of orderingFilter) Matches(object any) (bool, error) { return of.filter.Matches(object) }

func (of orderingFilter) SQLOrderBy() string { return of.order.SQLOrderBy() }

func (of orderingFilter) Less(object1, object2 any) (bool, error) {
    return of.order.Less(object1, object2)
}
