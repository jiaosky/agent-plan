package main

import (
	"fmt"
	"sort"
)

func main() {
	//nums := []int{-4, -1, -1, 0, 1, 2}
	//nums := []int{1, 1, 1}
	nums := []int{0, 0, 0, 0, 0, 1, 1, 2}
	fmt.Println(threeSum(nums))
}

// 借鉴 双数之和，变成三数和，先确定一个，然后倒排
func threeSum(nums []int) [][]int {
	return threeSumGeneral(nums, 0)
}

func threeSumGeneral(nums []int, target int) [][]int {

	var res [][]int
	sort.Ints(nums)

	i := 0
	for i < len(nums)-1 {
		first := nums[i]
		innerTarget := target - first

		left, right := i+1, len(nums)-1
		for left < right {
			l, r := left, right
			sum := nums[left] + nums[right]
			if sum == innerTarget {
				res = append(res, []int{nums[left], nums[right], nums[i]})
				for left < r && nums[left] == nums[l] {
					left++
				}
				for l < right && nums[right] == nums[r] {
					right--
				}
			} else if sum < innerTarget {
				for left < r && nums[left] == nums[l] {
					left++
				}
			} else {
				for l < right && nums[right] == nums[r] {
					right--
				}
			}
		}

		for i < len(nums)-1 && nums[i] == first {
			i++
		}
	}

	return res
}
