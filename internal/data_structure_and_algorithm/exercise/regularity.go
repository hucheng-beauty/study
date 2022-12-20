package exercise

import (
	"sort"
)

// SetZeroes 零矩阵
// 思路: 使用俩个数据标记矩阵中是否为零,然后遍历矩阵,0所在的行和列都置为0
func SetZeroes(matrix [][]int) {
	row := make([]bool, len(matrix))
	col := make([]bool, len(matrix[0]))
	for i, r := range matrix {
		for j, v := range r {
			if v == 0 {
				row[i] = true
				col[j] = true
			}
		}
	}
	for i, r := range matrix {
		for j := range r {
			if row[i] || col[j] {
				r[j] = 0
			}
		}
	}
}

func setZeroes(matrix [][]int) {
	n, m := len(matrix), len(matrix[0])
	row0, col0 := false, false

	// 标记第一行/列为0
	for _, v := range matrix[0] {
		if v == 0 {
			row0 = true
			break
		}
	}
	for _, r := range matrix {
		if r[0] == 0 {
			col0 = true
			break
		}
	}

	// 处理除了第一行之外的其他数据
	for i := 1; i < n; i++ {
		for j := 1; j < m; j++ {
			if matrix[i][j] == 0 {
				matrix[i][0] = 0
				matrix[0][j] = 0
			}
		}
	}
	for i := 1; i < n; i++ {
		for j := 1; j < m; j++ {
			if matrix[i][0] == 0 || matrix[0][j] == 0 {
				matrix[i][j] = 0
			}
		}
	}

	// 处理矩阵第一行或者第一列为 0
	if row0 {
		for j := 0; j < m; j++ {
			matrix[0][j] = 0
		}
	}
	if col0 {
		for _, r := range matrix {
			r[0] = 0
		}
	}
}

// IsStraight 是否为顺子,随机从扑克牌中抽五张牌
// 思路: 排序+遍历， 牌无重复， 时间O(nlogn), 空间O(1)
func IsStraight(nums []int) bool {
	//
	joker := 0
	sort.Ints(nums)
	for i := 0; i < 4; i++ {
		if nums[i] == 0 {
			joker++
		} else if nums[i] == nums[i+1] {
			return false
		}
	}
	return nums[4]-nums[joker] < 5
}

func divingBoard(shorter int, longer int, k int) []int {
	if k == 0 {
		return []int{}
	}
	if shorter == longer {
		return []int{shorter * k}
	}
	lengths := make([]int, k+1)
	for i := 0; i <= k; i++ {
		lengths[i] = shorter*(k-i) + longer*i
	}
	return lengths
}

func main() {
	/*matrix := make([][]int, 3)
	matrix[0] = []int{1, 2, 0}
	matrix[1] = []int{2, 3, 4}
	matrix[2] = []int{3, 4, 5}
	fmt.Println(matrix)
	SetZeroes(matrix)
	fmt.Println(matrix)*/

	/*nums := make([]int, 0)
	//nums = append(nums, 0, 2, 5, 1, 3)
	nums = append(nums, 0, 0, 5, 1, 3)
	fmt.Println(IsStraight(nums))*/
}
