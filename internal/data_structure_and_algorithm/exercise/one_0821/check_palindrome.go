package one_0821

/*
   给你一个字符串 s，最多 可以从中删除一个字符。
   请你判断 s 是否能成为回文字符串：如果能，返回 true ；否则，返回 false 。

   示例 1：
   输入：s = "aba"
   输出：true

   示例 2：
   输入：s = "abca"
   输出：true
   解释：你可以删除字符 'c' 。

   示例 3：
   输入：s = "abc"
   输出：false
*/

func isPalindrome1(str string, left, right int) bool {
	for left < right {
		if str[left] != str[right] {
			return false
		}
		left++
		right--
	}
	return true
}

func checkPalindrome(str string) bool {
	left, right := 0, len(str)-1
	for left < right {
		if str[left] != str[right] {
			// 判断是否为回文
			return isPalindrome1(str, left+1, right) ||
				isPalindrome1(str, left, right-1)
		}
		left++
		right--
	}
	return true
}
