package list

import "fmt"

// 找到链表的倒数第K个元素。

type List struct {
    Data int
    Next *List
}

// 1  2  3
// l1
// l2
func findKey(list *List, k int) *List {
    if list == nil || k <= 0 {
        return nil
    }

    l1 := list
    l2 := list

    for i := 0; i < k; i++ {
        if l1 == nil {
            return nil
        }
        l1 = l1.Next
    }
    for l1 != nil {
        l1 = l1.Next
        l2 = l2.Next
    }

    return l2
}

func main() {
    list := &List{
        Data: 1,
        Next: &List{
            Data: 2,
            Next: &List{
                Data: 3,
                Next: nil,
            },
        },
    }

    l := findKey(list, 5)

    fmt.Printf("%#v", l)
}
