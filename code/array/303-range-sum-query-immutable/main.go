package main

import "fmt"

func main() {
	// nums[0] = 0
	// nums[1] = -2
	// nums[2] = -2 + 0
	// nums[3] = -2 + 0 + 3
	// nums[4] = -2 + 0 + 3 -5
	// nums[5] = -2 + 0 + 3 -5 + 2
	// nums[6] = -2 + 0 + 3 -5 + 2 -1
	// left 2 ~ right 5  = 3 - 5 + 2 -1
	// num[6] - num[2]
	// [-2 + 0 + 3 -5 + 2 -1] - [-2 + 0]
	// right 的前序和  - [left-1] 的前序和
	// 0 ~ right 的前序和 = num[right-1]
	na := Constructor([]int{-2, 0, 3, -5, 2, -1})
	fmt.Println(na.SumRange(0, 2))
	fmt.Println(na.SumRange(2, 5))
	fmt.Println(na.SumRange(0, 5))
}

type NumArray struct {
	nums []int
}

func Constructor(nums []int) NumArray {
	// 定义一个数组，a[x] , i 是 a[0,..i-1] 的和
	// 如果 left 是 4 ，right 是 7
	// 就是 a[7] - a[4]
	a := make([]int, len(nums)+1)
	for i := 0; i < len(nums); i++ {
		a[i+1] = a[i] + nums[i]
	}

	return NumArray{nums: a}
}

func (this *NumArray) SumRange(left int, right int) int {
	// 定义一个数组，a[x] , i 是 a[0,..i-1] 的和
	// 如果 left 是 4 ，right 是 7
	// 就是 a[7] - a[4]
	return this.nums[right+1] - this.nums[left]
}
