package main

import "fmt"

/*
*
给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。

你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。

你可以按任意顺序返回答案。
*/
func twoSum(nums []int, target int) []int {

	result := make([]int, 2)
	m := make(map[int]int)

	for i := 0; i < len(nums); i++ {
		if k, exists := m[target-nums[i]]; exists {
			result[0] = i
			result[1] = k
			return result
		}
		m[nums[i]] = i
	}
	return nil
}

func twoSum2(nums []int, target int) []int {

	result := make([]int, 2)

	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				result[0] = i
				result[1] = j
				return result
			}
		}

	}
	return nil
}

func main() {
	var nums = []int{3, 2, 4}
	result := twoSum(nums, 6)
	fmt.Println(result)
}
