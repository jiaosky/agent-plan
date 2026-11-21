package main

import "fmt"

func main() {
	fmt.Println(containsNearbyDuplicate([]int{1, 2, 3, 1}, 3))
	fmt.Println(containsNearbyDuplicate([]int{1, 0, 1, 1}, 1))
	fmt.Println(containsNearbyDuplicate([]int{1, 2, 3, 1, 2, 3}, 2))
	fmt.Println(containsNearbyDuplicate([]int{1, 2, 3, 1, 1, 2, 3}, 2))
	fmt.Println(containsNearbyDuplicate([]int{1, 1}, 2))
}

func containsNearbyDuplicate(nums []int, k int) bool {
	// 每一个元素的角标
	find := make(map[int]int)
	left, right := 0, 0

	for right < len(nums) {
		c := nums[right]
		// 和前一个坐标相减
		if v, ok := find[c]; ok && right-v <= k {
			return true
		}
		find[c] = right
		for right-left >= k {
			d := nums[left]
			delete(find, d)
			left++
		}
		right++
	}

	return false
}
