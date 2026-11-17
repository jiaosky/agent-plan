package main

import "fmt"

func main() {
	num := []int{1, 1, 1, 2, 2, 3, 3, 5, 5, 5, 5, 5, 6, 6, 8, 8, 8, 9}
	result := removeDuplicates(num)
	fmt.Println(result)
	fmt.Println(num)
}

func removeDuplicates(nums []int) int {

	count := 1
	slow, fast := 0, 1

	for fast < len(nums) {
		if nums[fast] == nums[slow] {
			count++
			if count == 2 {
				slow++
				nums[slow] = nums[fast]
			}
		} else {
			slow++
			nums[slow] = nums[fast]
			count = 1
		}
		fast++
	}

	return slow + 1
}
