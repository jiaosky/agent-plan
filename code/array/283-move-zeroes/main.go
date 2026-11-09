package main

import "fmt"

func main() {

	num := []int{1, 2, 3, 1}
	moveZeroes(num)
	fmt.Println(num)
}

func moveZeroes(nums []int) {

	slow, fast := 0, 0

	for fast < len(nums) {
		if nums[fast] != 0 {
			nums[slow], nums[fast] = nums[fast], nums[slow]
			slow += 1
		}
		fast += 1
	}
}
