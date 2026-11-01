package main

import "fmt"

func main() {
	num := []int{-1, 1, 0, -3, 3}
	fmt.Println(productExceptSelf(num))
}

// 给你一个整数数组 nums，返回 数组 answer ，其中 answer[i] 等于 nums 中除了 nums[i] 之外其余各元素的乘积 。
// 题目数据 保证 数组 nums之中任意元素的全部前缀元素和后缀的乘积都在  32 位 整数范围内。
// 请 不要使用除法，且在 O(n) 时间复杂度内完成此题。
func productExceptSelf(nums []int) []int {

	left := make([]int, len(nums)+1)
	right := make([]int, len(nums)+1)

	left[0] = 1
	right[0] = 1

	for i := 0; i < len(nums); i++ {
		left[i+1] = left[i] * nums[i]
		right[i+1] = right[i] * nums[len(nums)-i-1]
	}

	fmt.Println(left)
	fmt.Println(right)

	ans := make([]int, len(nums))
	for i := 0; i < len(ans); i++ {
		ans[i] = left[i] * right[len(nums)-1-i]
	}

	return ans
}
