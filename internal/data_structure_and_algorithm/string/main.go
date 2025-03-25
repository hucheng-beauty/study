package main

import "fmt"

/*
   给你一个字符串 s，找到 s 中最长的回文子串。
   示例 1：
       输入：s = "babad"
       输出："bab"
       解释："aba" 同样是符合题意的答案。
   示例 2：
       输入：s = "cbbd"
       输出："bb"

*/

// babad
// l
// r

func findStr(s string) string {
    if len(s) <= 0 {
        return ""
    }

    maxLen := 1
    start := 0

    for i := 0; i < len(s); i++ {
        left, right := i, i
        for left > 0 && right < len(s) && s[left] == s[right] {
            if right-left+1 > maxLen {
                start = left
                maxLen = right - left + 1
            }
            left--
            right++
        }

        left, right = i, i+1
        for left >= 0 && right < len(s) && s[left] == s[right] {
            start = left
            maxLen = right - left + 1
        }
    }
    return s[start : start+maxLen]
}

func main() {
    s := "babad"
    fmt.Println(findStr(s)) // "bab"
}
