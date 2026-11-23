package main

import "fmt"

func main() {

	fmt.Println(minOperations([]int{1, 1, 4, 2, 3}, 5))
	fmt.Println(minOperations([]int{5, 6, 7, 8, 9}, 4))
	fmt.Println(minOperations([]int{3, 2, 20, 1, 1, 3}, 10))
	fmt.Println(minOperations([]int{5, 2, 3, 1, 1}, 5))
}

func minOperations(nums []int, x int) int {
	// 这道题看似是在最少次数可以获取目标值的和
	// 可以转化为 nums 的最大子数组 是 总sum - x
	// 这样就可以变成滑动窗口求和了

	sum := 0
	for _, n := range nums {
		sum += n
	}

	// 子数组中需要满足的最大值
	target := sum - x

	left, right := 0, 0
	maxLength := -1
	subSum := 0

	for right < len(nums) {
		c := nums[right]
		subSum += c
		right++

		// 如果子数组和大于 target ，就要缩小窗口，目标是 subSum == target 才对
		for left < right && subSum > target {
			d := nums[left]
			left++
			subSum -= d
		}

		//
		if target == subSum {
			if maxLength == -1 || maxLength < right-left {
				maxLength = right - left
			}
		}
	}

	if maxLength == -1 {
		return -1
	}

	return len(nums) - maxLength
}

func minOperations2(nums []int, x int) int {

	n := len(nums)
	sum := 0

	for _, num := range nums {
		sum += num
	}

	// 全部数组和 - 要移除元素的和
	target := sum - x

	left, right := 0, 0

	// 窗口内元素的和？
	windowSum := 0

	// 目标子数组的最大长度？
	maxLen := -1

	for right < n {
		// 要移除元素和相加
		windowSum += nums[right]
		right++

		// 已经移除元素的和 与 目标要移除的比较
		// 如果移出的多了，要还回来？
		for windowSum > target && left < right {
			windowSum -= nums[left]
			left++
		}

		if windowSum == target {
			if maxLen == -1 || right-left > maxLen {
				maxLen = right - left
			}
		}
	}

	if maxLen == -1 {
		return -1
	}

	// 字符串长度 - 统计出的最大子串
	return n - maxLen

}
