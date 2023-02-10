package list

import "fmt"

/*
    合并两个有序链表
    将两个升序链表合并为一个新的 升序 链表并返回。
    新链表是通过拼接给定的两个链表的所有节点组成的。
   示例 1：
       输入：l1 = [1,2,4], l2 = [1,3,4]
       输出：[1,1,2,3,4,4]
   示例 2：
       输入：l1 = [], l2 = []
       输出：[]
   示例 3：
       输入：l1 = [], l2 = [0]
       输出：[0]
*/

type Node struct {
    Val  int
    Next *Node
}

func MergeList(l1 *Node, l2 *Node) *Node {

    head := &Node{}
    index := head
    for l1 != nil && l2 != nil {
        if l1.Val < l2.Val {
            index.Next = l1
            l1 = l1.Next
        } else {
            index.Next = l2
            l2 = l2.Next
        }
        index = index.Next
    }
    if l1 != nil {
        index.Next = l1
    } else {
        index.Next = l2
    }
    return head
}

func TestMergeList() {
    l1 := &Node{
        Val: 1,
        Next: &Node{
            Val: 2,
            Next: &Node{
                Val:  4,
                Next: nil,
            },
        },
    }
    l2 := &Node{
        Val: 1,
        Next: &Node{
            Val: 3,
            Next: &Node{
                Val:  4,
                Next: nil,
            },
        },
    }

    list := MergeList(l1, l2)
    index := list

    for index.Next != nil {
        fmt.Println(index.Next.Val)
        index = index.Next
    }
}
