package main

import "fmt"

func main() {
	matrix := [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}
	fmt.Println(spiralOrder(matrix))
}

func spiralOrder(matrix [][]int) []int {

	// 横向的坐标
	m := len(matrix)
	n := len(matrix[0])

	// i
	upper := 0
	low := m - 1
	// j
	left := 0
	right := n - 1

	var res []int

	for len(res) < m*n {
		// 先横向
		if upper <= low {
			for j := left; j <= right; j++ {
				res = append(res, matrix[upper][j])
			}
			upper = upper + 1
		}

		if right >= left {
			for i := upper; i <= low; i++ {
				res = append(res, matrix[i][right])
			}
			right = right - 1
		}

		if low >= upper {
			for j := right; j >= left; j-- {
				res = append(res, matrix[low][j])
			}
			low = low - 1
		}

		if left <= right {
			for i := low; i >= upper; i-- {
				res = append(res, matrix[i][left])
			}
			left = left + 1
		}

	}
	return res
}
