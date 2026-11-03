package main

import "fmt"

func main() {
	length := 5
	update := [][]int{{1, 3, 2}, {2, 4, 3}, {0, 2, -2}}
	fmt.Println(getModifiedArray(length, update))
	// [-2,0,3,5,3]
}

// 有一个长度为 length 的数组 arr，初始时所有元素均为 0。
//
// 给你操作数组 updates，每个操作 [start, end, increment] 表示把闭区间 [start, end] 内的每个元素都加上 increment。按顺序完成所有操作后，返回数组 arr。

// -- core code template start
func getModifiedArray(length int, updates [][]int) []int {
	// 原数组 都是 0
	arr := make([]int, length)
	// 差分数组也都是 0
	diff := make([]int, length)

	for i := range updates {
		start := updates[i][0]
		end := updates[i][1]
		d := updates[i][2]

		diff[start] = diff[start] + d
		if end+1 < length {
			diff[end+1] = diff[end+1] - d
		}

	}

	arr[0] = diff[0]
	for i := 1; i < length; i++ {
		arr[i] = arr[i-1] + diff[i]
	}

	return arr
}
