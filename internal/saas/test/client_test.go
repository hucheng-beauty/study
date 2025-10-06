package test

import (
    "fmt"
    "testing"

    "github.com/samber/lo"
)

func TestForEach(t *testing.T) {
    m := map[string]string{
        "h1": "w1",
        "h2": "w2",
        "h3": "w3",
        "h4": "w4",
    }

    keys := lo.Keys(m)
    values := lo.Values(m)

    fmt.Println(keys)
    fmt.Println(values)
    fmt.Println(lo.Uniq(values))
}
