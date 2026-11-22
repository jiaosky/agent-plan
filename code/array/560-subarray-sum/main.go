package main

import "fmt"

func main() {
	nums := []int{1, -1, 0}
	k := 0
	fmt.Println(subarraySum(nums, k))
}

func subarraySum(nums []int, k int) int {
	count := 0
	preSum := preSumArray(nums)

	m := make(map[int]int, len(preSum))
	m[0] = 1
	for i := 1; i < len(preSum); i++ {

		need := preSum[i] - k
		if _, ok := m[need]; ok {
			count += m[need]
		}

		// 放入前缀和
		m[preSum[i]]++
	}

	return count

}

func preSumArray(nums []int) []int {
	preSum := make([]int, len(nums)+1)
	for i := 0; i < len(nums); i++ {
		preSum[i+1] = preSum[i] + nums[i]
	}
	return preSum
}

func subarraySum2(nums []int, k int) int {

	prsSum := preSumArray(nums)

	count := 0
	m := map[int]int{0: 1}
	for i := 1; i < len(prsSum); i++ {
		need := prsSum[i] - k
		if v, exist := m[need]; exist {
			count = count + v
		}

		if _, exist := m[prsSum[i]]; exist {
			m[prsSum[i]] = m[prsSum[i]] + 1
		} else {
			m[prsSum[i]] = 1
		}
	}

	return count

}

func preSumArray2(nums []int) []int {
	preSum := make([]int, len(nums)+1)
	preSum[0] = 0
	for i := range nums {
		preSum[i+1] = preSum[i] + nums[i]
	}
	return preSum
}
