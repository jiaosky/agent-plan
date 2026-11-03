package main

import (
	"fmt"
)

func main() {
	trips := [][]int{{9, 0, 1}, {3, 3, 7}}
	capacity := 4
	fmt.Println(carPooling(trips, capacity))
}

// 车上最初有 capacity 个空座位。车 只能 向一个方向行驶（也就是说，不允许掉头或改变方向）
// 给定整数 capacity 和一个数组 trips ,  trips[i] = [numPassengersi, fromi, toi] 表示第 i 次旅行有 numPassengersi 乘客，接他们和放他们的位置分别是 fromi 和 toi 。这些位置是从汽车的初始位置向东的公里数。
// 乘客在接上新乘客之前会被放下。在路线的每个点，车内乘客总数不得超过 capacity。
// 当且仅当你可以在所有给定的行程中接送所有乘客时，返回 true，否则请返回 false。
func carPooling(trips [][]int, capacity int) bool {

	// 构建差分数列
	// to 长度最大 1000 ，后面看看怎么优化？
	diff := make([]int, 1000)
	for _, trip := range trips {
		n := trip[0]
		from := trip[1]
		to := trip[2]

		diff[from] = diff[from] + n
		// 第 trip[2] 站乘客已经下车，
		// 即乘客在车上的区间是 [trip[1], trip[2] - 1]
		// 下车车站的后一站
		if (to-1)+1 < len(diff) {
			diff[(to-1)+1] = diff[(to-1)+1] - n
		}

	}

	num := make([]int, len(diff))
	num[0] = diff[0]
	if num[0] > capacity {
		return false
	}
	for i := 1; i < len(diff); i++ {
		num[i] = num[i-1] + diff[i]
		if num[i] > capacity {
			return false
		}
	}
	return true
}
