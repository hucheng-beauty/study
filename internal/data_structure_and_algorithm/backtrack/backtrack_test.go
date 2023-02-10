package backtrack

import (
    "fmt"
    "testing"
)

func TestMainer(t *testing.T) {
    s := new(Solution)
    fmt.Println(s.Permute([]int{1, 2, 3}))
}
