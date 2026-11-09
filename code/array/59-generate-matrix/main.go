package main

import "fmt"

func main() {
	fmt.Println(generateMatrix(0))
}

func generateMatrix(n int) [][]int {
	matrix := make([][]int, n)

	for i := range n {
		matrix[i] = make([]int, n)
	}

	upper := 0
	lower := n - 1
	left := 0
	right := n - 1

	num := 1
	for num <= n*n {
		if upper <= lower {
			for j := left; j <= right; j++ {
				matrix[upper][j] = num
				num += 1
			}
			upper += 1
		}

		if right >= left {
			for i := upper; i <= lower; i++ {
				matrix[i][right] = num
				num += 1
			}
			right -= 1
		}

		if lower >= upper {
			for j := right; j >= left; j-- {
				matrix[lower][j] = num
				num += 1
			}
			lower -= 1
		}

		if left <= right {
			for i := lower; i >= upper; i-- {
				matrix[i][left] = num
				num += 1
			}
			left += 1
		}
	}

	return matrix
}
