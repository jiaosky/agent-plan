package main

import "fmt"

func main() {
	weight := []int{1, 2, 3, 1, 1}
	days := 4

	//fmt.Println(f(weight, 1))

	fmt.Println(shipWithinDays(weight, days))
}

func f(weights []int, x int) int {
	// 函数应该怎么写呢？
	// 给定承载，计算天数
	// 包裹必须依次，并且不能拆开

	days := 0
	for i := 0; i < len(weights); {
		cap := x
		for i < len(weights) {
			// 超过运载边界
			if cap < weights[i] {
				break
			} else {
				cap -= weights[i]
			}
			i++
		}
		days++
	}

	return days
}

func shipWithinDays(weights []int, days int) int {
	// 求出在给点天数下，恰好可以运送完货物，每天最低的承载量
	// 每个包裹最小是 1 ，最大是 500

	// 承载量是变量，因此要根据承载量进行函数计算，计算是否等于给定的天数
	// 承载量越小，天数越大
	// 题目找到每天最低的承载量，因此就是二分查找，找到左边界
	left := 0
	// 注意，right 是开区间，所以额外加一
	right := 1

	for _, w := range weights {
		if left < w {
			left = w
		}
		right += w
	}

	// [0,501) 左闭右开区间

	for left < right {
		mid := left + (right-left)/2
		day := f(weights, mid)
		// 用时多，载重小了，寻找更大的载重
		if day == days {
			// 相等再找小的
			right = mid
		} else if day < days {
			right = mid
		} else if day > days {
			left = mid + 1
		}
	}

	return left
}
