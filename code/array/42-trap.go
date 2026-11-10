package main

import "fmt"

func main() {
	height := []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}
	fmt.Println(trap(height))
}

// 找到 i 左右各最大的高度，缓存起来
// min(l,r) - height[i] 就是当前格能够存水的数量
func trap(height []int) int {

	n := len(height)
	lm := make([]int, n)
	rm := make([]int, n)

	lm[0] = height[0]
	rm[n-1] = height[n-1]

	// 计算每一个 i 左右的最大值
	for i := 1; i < n; i++ {
		lm[i] = maxI(lm[i-1], height[i])
	}

	for i := n - 2; i >= 0; i-- {
		rm[i] = maxI(rm[i+1], height[i])
	}

	res := 0
	for i := 1; i < n; i++ {
		res += minI(lm[i], rm[i]) - height[i]
	}
	return res
}

func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minI(a, b int) int {
	if a < b {
		return a
	}
	return b
}
