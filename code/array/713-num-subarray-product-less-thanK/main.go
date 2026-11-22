package main

import "fmt"

func main() {
	//nums := []int{10, 5, 2, 6}
	//k := 100
	//
	nums := []int{686, 28, 455, 675, 605, 29, 942, 48, 502, 889, 854, 206, 231, 796, 272, 565, 887, 969, 558, 13, 22, 455, 145, 804, 15}
	k := 515854

	fmt.Println(numSubarrayProductLessThanK(nums, k))
}

func numSubarrayProductLessThanK(nums []int, k int) int {

	// 依旧是滑动数组
	// 进入窗口条件：乘积 < K
	// 退出窗口条件：乘积 >= K
	// 如果一个区间内的乘积 < K , 那么连续任意组合都小于K

	left, right := 0, 0
	result := 1
	count := 0

	for right < len(nums) {
		result = result * nums[right]
		right++

		for left < right && result >= k {
			result = result / nums[left]
			left++
		}

		// 此时窗口内肯定都是满足的
		// 这个是题目的关键：统计以右端点为结尾的子数组
		// 这样可以避免左侧端点重复
		count += right - left
	}

	return count
}
