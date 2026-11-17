package main

import "fmt"

func main() {
	num := []int{1, 0, 1, 2, 0, 0, 1, 1, 2, 0, 2, 1, 1, 0, 1, 2, 1, 1, 2, 0, 1, 1, 1}
	sortColors(num)
	fmt.Println(num)
}

// 红 白 蓝
// 0  1  2
func sortColors(nums []int) {
	left, right := 0, len(nums)-1

	for left < right {
		if nums[left] == 0 {
			left++
		} else if nums[left] == 2 {
			// 找到右边第一不为 2 的数
			for left < right && nums[right] == nums[left] {
				right--
			}
			// 仍旧是左小于右
			if left < right {
				nums[left], nums[right] = nums[right], nums[left]
				right--
			}
		} else {
			// 如果为 1 ，就找第一个为 0 的调换
			tmp := left + 1
			for tmp <= right {
				if nums[tmp] == 0 {
					nums[left], nums[tmp] = nums[tmp], nums[left]
					break
				}
				tmp++
			}
			left++
		}
	}
}
