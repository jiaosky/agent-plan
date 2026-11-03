package main

import "fmt"

func main() {
	matrix := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	rotate(matrix)
	fmt.Println(matrix)
}

// 给定一个 n × n 的二维矩阵 matrix 表示一个图像。请你将图像顺时针旋转 90 度。
//你必须在 原地 旋转图像，这意味着你需要直接修改输入的二维矩阵。请不要 使用另一个矩阵来旋转图像。

func rotate(matrix [][]int) {

	// 先对角镜像，然后每列进行互换
	for i := 0; i < len(matrix); i++ {
		for j := i + 1; j < len(matrix[i]); j++ {
			if i != j {
				matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
			}
		}
	}

	//printM(matrix)

	// 每列互换
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i])/2; j++ {
			matrix[i][j], matrix[i][len(matrix[i])-j-1] = matrix[i][len(matrix[i])-j-1], matrix[i][j]
		}
	}
	//printM(matrix)
}

func printM(matrix [][]int) {
	for _, v := range matrix {
		fmt.Println(v)
	}
	fmt.Println("--------------")
}
