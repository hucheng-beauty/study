package backtrack

/*
   回溯核心思想
        枚举所有可能的解;多阶段

*/

type Solution struct {
    result [][]int
}

// Permute 全排列
func (s *Solution) Permute(nums []int) [][]int {
    path := []int{}
    s.backtrack(nums, 0, path)
    return s.result
}

func (s *Solution) backtrack(nums []int /*可选列表,去除掉 path 中的数据*/, k int /*决策阶段*/, path []int /*记录路径*/) {
    // 结束条件
    if len(nums) == k {
        s.result = append(s.result, append([]int{}, path...))
        return
    }

    for i := 0; i < len(nums); i++ {
        if s.isExist(nums[i], path) {
            continue
        }
        // 做选择
        path = append(path, nums[i])
        // 递归
        s.backtrack(nums, k+1, path)
        // 撤销选择
        path = append(path[:len(path)-1], path[len(path)-1+1:]...)
    }
    return
}

func (s *Solution) isExist(key int, path []int) bool {
    for _, v := range path {
        if v == key {
            return true
        }
    }
    return false
}
