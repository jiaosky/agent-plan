package main

import "fmt"

func main() {

	fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))
	fmt.Println(twoSum([]int{2, 3, 4}, 6))
	fmt.Println(twoSum([]int{-1, 0}, -1))
	fmt.Println(twoSum([]int{5, 25, 75}, 100))
	fmt.Println(twoSum([]int{-1000, -1, 0, 1}, 1))
	fmt.Println(twoSum([]int{-5, -3, 0, 2, 4, 6, 8}, 5))
}
func twoSum(numbers []int, target int) []int {

	left := 0
	right := len(numbers) - 1

	for left <= right {
		sum := numbers[left] + numbers[right]
		if sum == target {
			break
		}
		if sum > target {
			right -= 1
		} else if sum < target {
			left += 1
		}

	}

	return []int{left + 1, right + 1}
}
