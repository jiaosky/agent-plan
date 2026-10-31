package main

import "fmt"

// 给你一个 m x n 的矩阵 mat 和一个整数 k ，请你返回一个矩阵 answer ，其中每个 answer[i][j] 是所有满足下述条件的元素 mat[r][c] 的和：
//i - k <= r <= i + k,
//j - k <= c <= j + k 且
//(r, c) 在矩阵内。
//示例 1：
//
//输入：mat = [[1,2,3],[4,5,6],[7,8,9]], k = 1
//输出：[[12,21,16],[27,45,33],[24,39,28]]
//示例 2：
//
//输入：mat = [[1,2,3],[4,5,6],[7,8,9]], k = 2
//输出：[[45,45,45],[45,45,45],[45,45,45]]

func main() {
	matrix := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	k := 2
	blockSum := matrixBlockSum(matrix, k)
	fmt.Println(blockSum)
}
func matrixBlockSum(mat [][]int, k int) [][]int {

	// answer[i][j] 在经过 r c  限制后，就变成了 mat 的子矩形
	// answer 的每一个子矩形，就是一个前缀和，answer 无非就是把每个前序和的结果，放在的一个矩形中
	// 因此还是要先算出来一个前缀和的矩阵
	ps := preSum(mat)

	// 构建 answer 矩阵
	ans := make([][]int, len(mat))
	for i := range ans {
		ans[i] = make([]int, len(mat[0]))
	}

	r := len(mat) - 1
	c := len(mat[0]) - 1

	for i := range ans {
		for j := range ans[i] {
			x1, y1 := i-k, j-k
			if x1 < 0 {
				x1 = 0
			}
			if y1 < 0 {
				y1 = 0
			}
			x2, y2 := i+k, j+k
			if x2 > r {
				x2 = r
			}
			if y2 > c {
				y2 = c
			}

			ans[i][j] = ps[x2+1][y2+1] - ps[x1][y2+1] - ps[x2+1][y1] + ps[x1][y1]
		}
	}
	return ans
}

// 构建前缀和矩阵
func preSum(mat [][]int) [][]int {
	row := len(mat)
	col := len(mat[0])

	a := make([][]int, row+1)
	for i := range a {
		a[i] = make([]int, col+1)
	}

	for i := 1; i <= row; i++ {
		for j := 1; j <= col; j++ {
			a[i][j] = mat[i-1][j-1] + a[i-1][j] + a[i][j-1] - a[i-1][j-1]
		}
	}

	return a
}
