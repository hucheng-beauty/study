package main

import "fmt"

// 二叉树的后续遍历
//  1
// 2 3
// 2 3 1

type Tree struct {
    Val   int
    Left  *Tree
    Right *Tree
}

func Print(t *Tree, r *[]int) []int {
    if t == nil {
        return nil
    }
    Print(t.Left, r)
    Print(t.Right, r)
    *r = append(*r, t.Val)
    return *r
}

func PrintTree(t *Tree) []int {
    if t == nil {
        return nil
    }
    r := []int{}
    s := []*Tree{}
    var p *Tree

    for t != nil || len(s) > 0 {
        for t != nil {
            s = append(s, t)
            t = t.Left
        }
        t = s[len(s)-1]
        s = s[:len(s)-1]

        if t.Right == nil || t.Right == p {
            r = append(r, t.Val)
            p = t
            t = nil
        } else {
            s = append(s, t)
            t = t.Right
        }
    }
    return r
}

func main() {
    r := []int{}
    t := &Tree{
        Val: 1,
        Left: &Tree{
            Val:   2,
            Left:  nil,
            Right: nil,
        },
        Right: &Tree{
            Val:   3,
            Left:  nil,
            Right: nil,
        },
    }
    // Print(t, r)
    fmt.Println(Print(t, &r))
    // fmt.Println(PrintTree(t))
}
