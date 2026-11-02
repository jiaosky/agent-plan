package main

import "fmt"

func main() {
	nums := []int{1, 1, 1}
	k := 2
	fmt.Println(subarraySum(nums, k))
}

func subarraySum(nums []int, k int) int {

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

func preSumArray(nums []int) []int {
	preSum := make([]int, len(nums)+1)
	preSum[0] = 0
	for i := range nums {
		preSum[i+1] = preSum[i] + nums[i]
	}
	return preSum
}
