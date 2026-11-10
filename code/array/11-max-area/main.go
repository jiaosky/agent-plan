package main

import "fmt"

func main() {
	height := []int{1, 8, 6, 2, 5, 4, 8, 3, 7}
	fmt.Println(maxArea(height))
}

func maxArea(height []int) int {

	if len(height) < 2 {
		return 0
	}

	lm := height[0]
	rm := height[len(height)-1]

	left := 0
	right := len(height) - 1
	area := 0
	for left < right {
		lm = max(lm, height[left])
		rm = max(rm, height[right])

		area = max(area, min(lm, rm)*(right-left))

		if lm < rm {
			left++
		} else {
			right--
		}

	}

	return area
}
