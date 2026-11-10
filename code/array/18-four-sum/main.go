package main

import (
	"fmt"
	"sort"
)

func main() {
	nums := []int{-2, -1, 0, 0, 1, 2}
	fmt.Println(fourSum(nums, 0))
}

func fourSum(nums []int, target int) [][]int {

	var res [][]int
	sort.Ints(nums)

	i := 0
	for i < len(nums)-1 {
		first := nums[i]
		innerTarget := target - first

		r := threeSumGeneral(nums, i+1, innerTarget)
		for _, v := range r {
			res = append(res, append(v, nums[i]))
		}

		for i < len(nums)-1 && nums[i] == first {
			i++
		}
	}

	return res
}

func threeSumGeneral(nums []int, i, target int) [][]int {
	var res [][]int
	for i < len(nums)-1 {
		first := nums[i]
		innerTarget := target - first

		rs := twoSumGeneral(nums, i+1, innerTarget)
		for _, v := range rs {
			res = append(res, append(v, nums[i]))
		}

		for i < len(nums)-1 && nums[i] == first {
			i++
		}
	}

	return res
}

func twoSumGeneral(nums []int, i, target int) [][]int {
	var res [][]int
	left, right := i, len(nums)-1
	for left < right {
		l, r := left, right
		sum := nums[left] + nums[right]
		if sum == target {
			res = append(res, []int{nums[left], nums[right]})
			for left < r && nums[left] == nums[l] {
				left++
			}
			for l < right && nums[right] == nums[r] {
				right--
			}
		} else if sum < target {
			for left < r && nums[left] == nums[l] {
				left++
			}
		} else {
			for l < right && nums[right] == nums[r] {
				right--
			}
		}
	}

	return res
}
