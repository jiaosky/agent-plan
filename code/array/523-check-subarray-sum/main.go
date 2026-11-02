package main

import "fmt"

func main() {
	nums := []int{23, 2, 4, 6, 6}
	k := 7
	fmt.Println(checkSubarraySum(nums, k))
	fmt.Println(checkSubarraySum2(nums, k))
}

func checkSubarraySum(nums []int, k int) bool {
	ps := preSumArray(nums)
	m := make(map[int]int)
	// 要从 0 开始，因为要考虑后面前缀和刚好是K倍数，此时需要余数 0 ！！
	for i := 0; i < len(ps); i++ {
		yushu := ps[i] % k
		if _, exist := m[yushu]; !exist {
			m[yushu] = i
		}
	}

	for i := 1; i < len(ps); i++ {
		// 计算 need，使得 (preSum[i] - need) % k == 0
		need := ps[i] % k
		if idx, exists := m[need]; exists {
			if i-idx >= 2 {
				// 这个子数组的长度至少为 2
				return true
			}
		}
	}

	return false
}

func checkSubarraySum2(nums []int, k int) bool {
	ps := preSumArray(nums)
	m := make(map[int]int)
	for i := 0; i < len(ps); i++ {
		yushu := ps[i] % k
		if _, exist := m[yushu]; !exist {
			m[yushu] = i
		}

		if idx, exists := m[yushu]; exists {
			if i-idx >= 2 {
				return true
			}
		}
	}

	//for i := 1; i < len(ps); i++ {
	//	yushu := ps[i] % k
	//
	//	if idx, exists := m[yushu]; exists {
	//		if i-idx >= 2 {
	//			return true
	//		}
	//	}
	//}
	return false
}

func preSumArray(nums []int) []int {
	preSum := make([]int, len(nums)+1)
	preSum[0] = 0
	for i := range nums {
		preSum[i+1] = preSum[i] + nums[i]
	}
	return preSum
}
