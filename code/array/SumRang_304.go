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

	//        matrix 的下标是从 0 开始：
	//
	//        c=0 c=1 c=2 c=3 c=4
	//r=0      3   0   1   4   2
	//r=1      5   6   3   2   1
	//r=2      1   2   0   1   5
	//r=3      4   1   0   1   7
	//r=4      1   0   3   0   5

	// 前缀和 preSum[i][j] 指的是角标从 matrix 左上角[0][0] 到 matrix 右下角[i][j] 全部元素的和
	// 也就是说 前缀和 的左边是右下角
	// 因此在计算给定目标矩形时：preSum[右下角] - preSum[左下] - preSum[右上] + preSum[左上角]
	// preSum[x2,y2]-preSum[x2,y1-1]-preSum[x1-1,y2]+preSum[x1-1,y1-1]
	// 对于 [0,0] ~ [0,0] 就是
	// preSum[0,0] - preSum[0,-1] - preSum[-1,0] - preSum[-1,-1]
	// 为了避免出现负数
	// preSum 就需要多一行列
	// preSum 的角标 = 原数组 -1
	// preSum[1,1] = preSum[0,0]

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
