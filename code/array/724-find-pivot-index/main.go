package main

import "fmt"

func main() {
	nums := []int{-1, -1, -1, 1, 1, 1}
	fmt.Println(pivotIndex(nums))
}

// 给你一个整数数组 nums ，请计算数组的 中心下标 。
// 数组 中心下标 是数组的一个下标，其左侧所有元素相加的和等于右侧所有元素相加的和。
// 如果中心下标位于数组最左端，那么左侧数之和视为 0 ，因为在下标的左侧不存在元素。这一点对于中心下标位于数组最右端同样适用。
// 如果数组有多个中心下标，应该返回 最靠近左边 的那一个。如果数组不存在中心下标，返回 -1 。
func pivotIndex(nums []int) int {

	// 计算全部的前序和
	// 最后一个是全部元素的和
	preSum := make([]int, len(nums)+1)
	for i := range nums {
		preSum[i+1] = preSum[i] + nums[i]
	}

	maxSum := preSum[len(nums)]

	if 0 == maxSum-nums[0] {
		return 0
	}

	// 从左到右依次进行前序和和最大元素-左边移动的元素对比
	for i := 0; i < len(nums)-1; i++ {
		// 前序和  = 最大值 - 前序和 - 当前元素
		if preSum[i+1] == maxSum-preSum[i+2] {
			return i + 1
		}
	}
	return -1

}
