package main

import "fmt"

func main() {
	nums := []int{-7, -3, 2, 3, 11}
	fmt.Println(sortedSquares(nums))
}

func sortedSquares(nums []int) []int {
	result := make([]int, len(nums))

	left, right := 0, len(nums)-1
	index := len(result) - 1
	for left <= right {
		if nums[left]*nums[left] > nums[right]*nums[right] {
			result[index] = nums[left] * nums[left]
			left++
		} else {
			result[index] = nums[right] * nums[right]
			right--
		}
		index--
	}
	return result
}
