package main

import "fmt"

func main() {
	num1 := []int{0}
	m := 0
	num2 := []int{1}
	n := 1
	merge(num1, m, num2, n)
	fmt.Println(num1)
}

// 从后向前放
func merge(nums1 []int, m int, nums2 []int, n int) {
	i := m - 1
	j := n - 1

	index := len(nums1) - 1

	for i >= 0 && j >= 0 {
		if nums2[j] > nums1[i] {
			nums1[index] = nums2[j]
			j--
		} else {
			nums1[index] = nums1[i]
			i--
		}

		index--
	}

	for j >= 0 {
		nums1[index] = nums2[j]
		j--
		index--

	}
}
