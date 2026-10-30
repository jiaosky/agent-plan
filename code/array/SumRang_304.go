package main

import "fmt"

func main() {
	na := ConstructorMatrix([][]int{{3, 0, 1, 4, 2}, {5, 6, 3, 2, 1}, {1, 2, 0, 1, 5}, {4, 1, 0, 1, 7}, {1, 0, 3, 0, 5}})
	fmt.Println(na.SumRegion(2, 1, 4, 3))
	//fmt.Println(na.SumRegion(1, 1, 2, 2))
	//fmt.Println(na.SumRegion(1, 2, 2, 4))
}

type NumMatrix struct {
	nums [][]int
}

func ConstructorMatrix(matrix [][]int) NumMatrix {

	m := len(matrix)    // 数组的行数
	n := len(matrix[0]) // 数组的列数

	preSum := make([][]int, m+1)
	for i := range preSum {
		preSum[i] = make([]int, n+1)
	}

	// 这道题思路和 SunRange_303 是相似的
	// 都是算出角标[i][j] 前序和
	// 然后再计算 [i+1][j+1] 的时候，使用 matrix[i+1][j+1] + preSum[i][j+1] + preSum[i+1][j] - preSum[i][j] 计算出新元素的前序和
	// 恶心的地方在于，目标矩形的转化
	// preSum[1][1] 存放的是 matrix[0][0]的前序和
	// preSum[2][2] 存放的是 matrix[1][1]的前序和
	// preSum[2][3] 存放的是 matrix[1][2]的前序和
	// preSum 的角标，比 matrix 的角标多一个
	// 目标[1,2,2,4] 转化为原矩形是 matrix[1,2] matrix[2,4]
	// 但是映射到 preSum[2,3] 和 preSum[3,5]
	// 因此在 sumRegion 的时候，就需要 row 和 col 分别 + 1 去处理了
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			preSum[i][j] = matrix[i-1][j-1] + preSum[i-1][j] + preSum[i][j-1] - preSum[i-1][j-1]
		}
	}
	return NumMatrix{preSum}
}

func (this *NumMatrix) SumRegion(row1 int, col1 int, row2 int, col2 int) int {
	return this.nums[row2+1][col2+1] - this.nums[row1][col2+1] - this.nums[row2+1][col1] + this.nums[row1][col1]
}
