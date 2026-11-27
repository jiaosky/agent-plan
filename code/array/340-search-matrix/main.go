package main

import (
	"fmt"
)

func main() {
	//matrix := [][]int{{1, 4, 7, 11, 15}, {2, 5, 8, 12, 19}, {3, 6, 9, 16, 22}, {10, 13, 14, 17, 24}, {18, 21, 23, 26, 30}}
	//target := 100

	matrix := [][]int{{1, 4, 7, 11, 15}, {2, 5, 8, 12, 19}, {3, 6, 9, 16, 22}, {10, 13, 14, 17, 24}, {18, 21, 23, 26, 30}}
	target := 19

	fmt.Println(searchMatrix(matrix, target))

}

func searchMatrix(matrix [][]int, target int) bool {

	m, n := len(matrix[0]), len(matrix)

	// 从 [n,0] 开始，向上是递减，向右是递增

	i := n - 1
	j := 0

	for i >= 0 && j < m {
		if matrix[i][j] == target {
			return true
		} else if matrix[i][j] > target {
			// 右下角 > 目标值，向上找
			if upper(matrix, i, j, target) {
				return true
			}

		} else {
			// 右下角 < 目标值，向右找
			if goRight(matrix, i, j, target) {
				return true
			}
		}
		// 没有找到，进入下一个区间
		i--
		j++
	}

	return false
}

// 向🈶右的二分查找
func goRight(matrix [][]int, i, j int, target int) bool {
	left, right := j, len(matrix[0])

	for left < right {
		mid := left + (right-left)/2
		if matrix[i][mid] == target {
			return true
		} else if matrix[i][mid] > target {
			// 右下角 > 目标值，向上找
			right = mid
		} else {
			left = mid + 1
		}
	}
	return false
}

// 向上的二分查找
func upper(matrix [][]int, i, j int, target int) bool {
	left, right := 0, i+1

	for left < right {
		mid := left + (right-left)/2
		if matrix[mid][j] == target {
			return true
		} else if matrix[mid][j] > target {
			// 右下角 > 目标值，向上找
			right = mid
		} else {
			left = mid + 1
		}
	}
	return false
}
