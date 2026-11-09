package main

import (
	"fmt"
)

func main() {

	num := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	k := removeDuplicates(num)
	fmt.Println(num)
	fmt.Println(k)
}

func removeDuplicates(nums []int) int {
	index := 0 // 新数组的坐标
	j := 0

	m := map[int]int{}
	for j < len(nums) {
		if _, exist := m[nums[j]]; !exist {
			m[nums[j]] = j
			nums[index] = nums[j]
			index += 1
			j += 1
		} else {
			j += 1
		}
	}
	return index
}
