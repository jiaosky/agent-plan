package main

import "fmt"

func main() {
	//matrix := [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}
	//target := 3

	matrix := [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}
	target := 13
	fmt.Println(searchMatrix(matrix, target))
}

func searchMatrix(matrix [][]int, target int) bool {
	left, right := 0, len(matrix)*len(matrix[0])
	for left < right {
		mid := left + (right-left)/2

		i, j := mid/len(matrix[0]), mid%len(matrix[0])
		if target == matrix[i][j] {
			return true
		} else if target > matrix[i][j] {
			// 目标值大于中位值，应该放大左边界
			left = mid + 1
		} else {
			right = mid
		}
	}

	return false
}
