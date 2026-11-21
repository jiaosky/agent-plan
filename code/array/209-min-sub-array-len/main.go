package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(minSubArrayLen(7, []int{2, 3, 1, 2, 4, 3}))
	fmt.Println(minSubArrayLen(4, []int{1, 4, 4}))
	fmt.Println(minSubArrayLen(11, []int{1, 1, 1, 1, 1, 1, 1, 1}))
}

func minSubArrayLen(target int, nums []int) int {
	// 滑动窗口
	// sum 记录 left ~ right 的值
	// 每次 sum >= target 记录数组长度,找到最小的 sum
	// 循环条件就是 sum < target ,right ++

	left, right := 0, 0
	minLength := math.MaxInt
	sum := 0
	for right < len(nums) {
		sum += nums[right]
		right++
		// 此时窗口内的总和 sum >= target 了
		// 锁小窗口,当 sum 小于 target 是退出
		for sum >= target {
			if minLength > right-left {
				minLength = right - left
			}
			sum = sum - nums[left]
			left++
		}
	}

	if minLength == math.MaxInt {
		return 0
	}

	return minLength
}

func minSubArrayLen3(target int, nums []int) int {
	left, right := 0, 0
	minLength := math.MaxInt
	sum := 0

	for ; right < len(nums); right++ {
		sum += nums[right]
		for sum >= target {
			minLength = min(minLength, right-left+1)
			sum = sum - nums[left]
			left++
		}
	}

	if minLength == math.MaxInt {
		return 0
	}

	return minLength
}

func minSubArrayLen2(target int, nums []int) int {
	// 滑动窗口
	// sum 记录 left ~ right 的值
	// 每次 sum >= target 记录数组长度,找到最小的 sum
	// 循环条件就是 sum < target ,right ++

	left, right := 0, 0
	minLength := math.MaxInt
	sum := 0
	for right < len(nums) {
		for right < len(nums) && sum < target {
			c := nums[right]
			sum += c
			right++
		}
		// 此时窗口内的总和 sum >= target 了
		// 锁小窗口,当 sum 小于 target 是退出
		for sum >= target {
			if minLength > right-left {
				minLength = right - left
			}

			d := nums[left]
			sum = sum - d
			left++
		}
	}

	if minLength == math.MaxInt {
		return 0
	}

	return minLength
}
