package main

import "fmt"

func Counter() func() (int, int) {
    a := 0
    return func() (int, int) {
        b := 0
        a++
        b++
        return a, b
    }
}

func Print() {
    counter := Counter()
    for i := 0; i < 4; i++ {
        /*
           // 通过闭包修改 a 的值
              1,1
              2,1
              3,1
              4,1
        */
        fmt.Println(counter()) // 请回答打印结果
    }
}

func SliExpansion() {
    sliceA := []int{1, 2, 3, 4, 5}
    sliceB := sliceA
    sliceA[0] = -1
    fmt.Println(sliceB[0])

    sliceB = append(sliceB, 1)
    sliceB[0] = -2
    fmt.Println(sliceA[0])
}

func SliMake() {
    sli := make([]int, 3)
    sli = append(sli, 1)
    sli = append(sli, 2)
    sli = append(sli, 3)
    fmt.Println(sli)
}

func main() {
    Print()
    SliExpansion()
    SliMake()
}
