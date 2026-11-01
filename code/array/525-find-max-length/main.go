package main

import "fmt"

func main() {
	nums := []int{1, 0, 0, 1}
	fmt.Println(findMaxLength(nums))
}

// 给定一个二进制数组 nums , 找到含有相同数量的 0 和 1 的最长连续子数组，并返回该子数组的长度。示例 1：
//
// 输入：nums = [0,1]
// 输出：2
// 说明：[0, 1] 是具有相同数量 0 和 1 的最长连续子数组。
// 示例 2：
//
// 输入：nums = [0,1,0]
// 输出：2
// 说明：[0, 1] (或 [1, 0]) 是具有相同数量 0 和 1 的最长连续子数组。
// 示例 3：
//
// 输入：nums = [0,1,1,1,1,1,0,0,0]
// 输出：6
// 解释：[1,1,1,0,0,0] 是具有相同数量 0 和 1 的最长连续子数组。
func findMaxLength(nums []int) int {

	preSum := make([]int, len(nums)+1)

	for i := range nums {
		if nums[i] == 0 {
			preSum[i+1] = preSum[i] - 1
		} else {
			preSum[i+1] = preSum[i] + 1
		}
	}

	maxLength := 0
	m := make(map[int]int)
	for i := 0; i < len(preSum); i++ {
		if _, exist := m[preSum[i]]; !exist {
			m[preSum[i]] = i
		} else {
			length := i - m[preSum[i]]
			if length > maxLength {
				maxLength = length
			}
		}
	}

	return maxLength
}
