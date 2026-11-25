package main

import "fmt"

func main() {
	//nums := []int{7, 2, 5, 10, 8}
	//k := 2

	nums := []int{1, 2, 3, 4, 5}
	k := 2
	fmt.Println(splitArray(nums, k))
}

func splitArray(nums []int, k int) int {

	// 传入的 K 是结果，结果是切分的数组数量

	// f(nums,x) return k
	// x 是满足 k 份下，其中子数组和最大值是所有切分方法的最小值

	// 类比搬运货物，给定 K 就是搬运天数
	// x 就是搬运力的最小值
	// 函数是计算在 x 下，搬运完全部东西的天数
	// 找到 x 的最小值

	left := 0
	right := 1 // 左闭右开区间
	for v := range nums {
		// 要满足可以至少装走一个格子的东西
		if left < nums[v] {
			left = nums[v]
		}
		right += nums[v]
	}

	for left < right {
		mid := left + (right-left)/2
		ks := f(nums, mid)
		if ks == k {
			// 天数相同，尝试用更小的运载量实现
			right = mid
		} else if ks > k {
			// 天数大于目标，说明运载量太小了，需要放大运载量
			left = mid + 1
		} else if ks < k {
			// 天数小于目标，说明运载量太大了，需要减少运载量
			right = mid
		}
	}
	return left
}

func f(nums []int, x int) int {
	k := 0
	for i := 0; i < len(nums); {
		cap := x
		// 可以一波带走
		for i < len(nums) {
			if cap >= nums[i] {
				cap -= nums[i]
			} else {
				break
			}
			i++
		}
		k++
	}
	return k
}
