package recursion

/*
	怎么发现这个问题可以用递归来做?
		规模更小的问题,跟规模大点的问题,解决思路相同、仅仅规模不同
		利用子问题的解可以组合得到原问题的解
		存在最小子问题,可以直接返回结果,即存在递归终止条件
	存在问题:重复计算
	递归相关题型和解题思路
		重复结构
		递推公式
		递归终止条件
	相关例题
		爬楼梯
		细胞分裂
		逆序打印链表
*/

// Fibonacci O(n^2) C(n)
func Fibonacci(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	if n == 2 {
		return 2
	}

	return Fibonacci(n-1) + Fibonacci(n-2)
}
