package registry

import (
    "fmt"
    "reflect"
    "strings"
)

var errFieldNotFound = fmt.Errorf("field with expected name tag not found")

type in struct {
    tableBound    string
    fieldName     string
    allowedValues []any
}

// In selects entities where fieldName matches at least one of allowedValues.
// If len(allowedValues) = 0, this filter is equivalent to that All filter.
func In(tableBound, fieldName string, allowedValues ...any) Filter {
    return in{
        tableBound:    tableBound,
        fieldName:     fieldName,
        allowedValues: allowedValues,
    }
}

func (f in) SQL() (string, []any) {
    if f.allowedValues == nil {
        return All().SQL()
    }

    fieldName := fullFieldName(f.tableBound, f.fieldName)
    switch len(f.allowedValues) {
    case 0:
        return None().SQL()
    case 1:
        return fmt.Sprintf("%s = ?", fieldName), f.allowedValues
    default:

        placeholders := strings.Join(strings.Split(strings.Repeat("?", len(f.allowedValues)), ""), ", ")
        return fmt.Sprintf("%s IN (%s)", fieldName, placeholders), f.allowedValues
    }
}

// Matches checks if object contains a field with struct tag `db`
// which is equal to `fieldName` and value in `allowedValues` slice.
func (f in) Matches(obj any) (bool, error) {
    if f.allowedValues == nil {
        return All().Matches(obj)
    }
    if len(f.allowedValues) == 0 {
        return None().Matches(obj)
    }

    field, found, err := findFieldByName(f.fieldName, obj)
    if err != nil {
        return false, fmt.Errorf("find field by name: %w", err)
    }
    if !found {
        return false, fmt.Errorf("[type=%T, field_name=%v] %w", obj, f.fieldName, errFieldNotFound)
    }

    return f.acceptedValue(field)
}

func (f in) acceptedValue(v any) (bool, error) {
    for _, allowedValue := range f.allowedValues {
        // if type of v is *T and type of allowed value is T, dereference v and compare values.
        checkedValue := v
        if reflect.TypeOf(checkedValue) == reflect.PtrTo(reflect.TypeOf(allowedValue)) {
            if reflect.ValueOf(checkedValue).IsZero() {
                return false, nil
            }

            checkedValue = reflect.ValueOf(checkedValue).Elem().Interface()
        }
        allowedT, checkedT := reflect.TypeOf(allowedValue), reflect.TypeOf(checkedValue)
        if allowedT != checkedT {
            return false, fmt.Errorf("[allowed=%v, checked=%v] allowed and expected types don't match",
                allowedT, checkedT)
        }

        if !allowedT.Comparable() {
            return false, fmt.Errorf("[type=%T] type is not comparable", checkedValue)
        }
        if allowedValue == checkedValue {
            return true, nil
        }
    }

    return false, nil
}

func LessThan(tableBound, fieldName string, comparedValue any) Filter {
    return comparison(tableBound, fieldName, compareMethodLessThan{}, comparedValue)
}

func GreaterThan(tableBound, fieldName string, comparedValue any) Filter {
    return comparison(tableBound, fieldName, compareMethodGreaterThan{}, comparedValue)
}

// compareMethod describes a valid SQL comparison. It defines a Compare function
// that should return the same result as cmp in SQL.
type compareMethod interface {
    // SQL function defines what should be placed in query between two object to make it a valid cmp.
    // For instance, for 'less than' condition, valid condition is "o1 < o2" - in that case, function
    // should return "<". For 'less than or equal', it would be "<=", and so on.
    SQL() string

    // Compare checks the result of the cmp, using supplied arguments. It accepts only
    // numeric objects (integers, unsigned integers, float numbers).
    Compare(o1, o2 any) (bool, error)
}

var _ compareMethod = compareMethodLessThan{}

const (
    GreaterOperator = ">"
    LessOperator    = "<"
)

// compareMethodLessThan defines a "<" cmp condition.
type compareMethodLessThan struct{}

func (compareMethodLessThan) SQL() string {
    return LessOperator
}

func (compareMethodLessThan) Compare(o1, o2 any) (bool, error) {
    res, err := compareBuiltin(o1, o2)
    if err != nil {
        return false, err
    }
    return res == compareResultLessThan, nil
}

// compareMethodGreaterThan defines a ">" cmp condition.
type compareMethodGreaterThan struct{}

func (compareMethodGreaterThan) SQL() string {
    return GreaterOperator
}

func (compareMethodGreaterThan) Compare(o1, o2 any) (bool, error) {
    res, err := compareBuiltin(o1, o2)
    if err != nil {
        return false, err
    }
    return res == compareResultGreaterThan, nil
}

type cmp struct {
    tableBound    string
    fieldName     string
    compareMethod compareMethod
    comparedValue any
}

// comparison is a condition that, based on a supplied comparison method,
// compares the value of the field value passed on initialization.
func comparison(tableBound, fieldName string, compareMethod compareMethod, comparedValue any) Filter {
    return cmp{
        tableBound:    tableBound,
        fieldName:     fieldName,
        compareMethod: compareMethod,
        comparedValue: comparedValue,
    }
}

func (f cmp) SQL() (string, []any) {
    fieldName := fullFieldName(f.tableBound, f.fieldName)
    return fmt.Sprintf("%s %s ?", fieldName, f.compareMethod.SQL()), []any{f.comparedValue}
}

func (f cmp) Matches(obj any) (bool, error) {
    field, found, err := findFieldByName(f.fieldName, obj)
    if err != nil {
        return false, fmt.Errorf("find field by name: %w", err)
    }
    if !found {
        return false, fmt.Errorf("[type=%T, field_name=%v] %w", obj, f.fieldName, errFieldNotFound)
    }

    rt, err := f.compareMethod.Compare(field, f.comparedValue)
    if err != nil {
        return false, fmt.Errorf("compare method evaluate: %w", err)
    }

    return rt, nil
}

func fullFieldName(tableBound, fieldName string) string {
    rt := fieldName
    if tableBound != "" {
        rt = tableBound + "." + rt
    }

    return rt
}

func findFieldByName(fieldName string, object any) (any, bool, error) {
    objT, objV := reflect.TypeOf(object), reflect.ValueOf(object)

    if objT.Kind() != reflect.Struct {
        return nil, false, fmt.Errorf("[type=%T] only matching structs is supported", object)
    }

    for i := 0; i < objT.NumField(); i++ {
        fieldType := objT.Field(i)

        ti, ok, err := dbTagInfo(fieldType)
        if err != nil {
            return nil, false, fmt.Errorf("[type=%T, field_name=%v] parse db tag info: %w",
                object, objT.Field(i).Name, err)
        }

        if !ok {
            continue
        }

        if fieldType.Anonymous {
            field, found, errF := findFieldByName(fieldName, objV.Field(i).Interface())
            switch {
            case errF == nil:
                if !found {
                    continue
                }

                return field, true, nil

            default:
                return nil, false, fmt.Errorf("[field=%v] match inline field: %w",
                    objT.Field(i).Name, err)
            }
        }

        if ti.FieldName == fieldName {
            return objV.Field(i).Interface(), true, nil
        }
    }

    return nil, false, nil
}

const (
    compareResultLessThan    compareResult = -1
    compareResultEqual       compareResult = 0
    compareResultGreaterThan compareResult = 1
)

type compareResult int

func compareBuiltin(object1, object2 any) (compareResult, error) {
    switch v := object1.(type) {
    case string:
        v2, ok := object2.(string)
        if !ok {
            return 0, fmt.Errorf("[type=%T value=%v] object2 is not string", object2, object2)
        }

        return compareString(v, v2), nil
    case int, int8, int16, int32, int64:
        v1, err := toInt64(v)
        if err != nil {
            panic(err)
        }

        v2, err := toInt64(object2)
        if err != nil {
            return 0, fmt.Errorf("to int64: %w", err)
        }

        return compareInt64(v1, v2), nil
    case uint, uint8, uint16, uint32, uint64:
        v1, err := toUint64(v)
        if err != nil {
            panic(err)
        }

        v2, err := toUint64(object2)
        if err != nil {
            return 0, fmt.Errorf("to uint64: %w", err)
        }

        return compareUint64(v1, v2), nil
    case float32, float64:
        v1, err := toFloat64(v)
        if err != nil {
            panic(err)
        }

        v2, err := toFloat64(object2)
        if err != nil {
            return 0, fmt.Errorf("to float64: %w", err)
        }

        return compareFloat64(v1, v2), nil

    default:
        return 0, fmt.Errorf("[o1_type=%T] unsupported type for cmp", object1)
    }
}

func toInt64(x any) (int64, error) {
    switch n := x.(type) {
    case int:
        return int64(n), nil
    case int8:
        return int64(n), nil
    case int16:
        return int64(n), nil
    case int32:
        return int64(n), nil
    case int64:
        return n, nil
    }
    return 0, fmt.Errorf("[type=%T] cannot convert to int64", x)
}

func compareInt64(i1, i2 int64) compareResult {
    if i1 < i2 {
        return compareResultLessThan
    } else if i1 == i2 {
        return compareResultEqual
    }

    return compareResultGreaterThan
}

func toUint64(x any) (uint64, error) {
    switch n := x.(type) {
    case uint:
        return uint64(n), nil
    case uintptr:
        return uint64(n), nil
    case uint8:
        return uint64(n), nil
    case uint16:
        return uint64(n), nil
    case uint32:
        return uint64(n), nil
    case uint64:
        return n, nil
    }

    return 0, fmt.Errorf("[type=%T] cannot convert to uint64", x)
}

func compareUint64(u1, u2 uint64) compareResult {
    if u1 < u2 {
        return compareResultLessThan
    } else if u1 == u2 {
        return compareResultEqual
    }

    return compareResultGreaterThan
}

func toFloat64(x any) (float64, error) {
    switch n := x.(type) {
    case float32:
        return float64(n), nil
    case float64:
        return n, nil
    }

    return 0, fmt.Errorf("[type=%T] cannot convert to float64", x)
}

func compareFloat64(f1, f2 float64) compareResult {
    if f1 < f2 {
        return compareResultLessThan
    } else if f1 == f2 {
        return compareResultEqual
    }

    return compareResultGreaterThan

}

func compareString(s1, s2 string) compareResult {
    if s1 < s2 {
        return compareResultLessThan
    } else if s1 == s2 {
        return compareResultEqual
    }

    return compareResultGreaterThan
}
