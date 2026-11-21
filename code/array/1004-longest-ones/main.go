package main

import "fmt"

func main() {
	fmt.Println(longestOnes4([]int{1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0}, 2))
	//fmt.Println(longestOnes2([]int{0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1}, 3))
}

func longestOnes4(nums []int, k int) int {
	left := 0
	for _, x := range nums {
		k = k - 1 - x // 借用了几次，因为 0 会借用 1次，1 不借用
		if k < 0 {    // 额度不够时，窗口不缩小，而是整体向右平移 1 位
			k = k + 1 - nums[left]
			left++
		}
	}
	return len(nums) - left
}

func longestOnes3(nums []int, k int) int {
	left, right := 0, 0
	// 记录窗口中 1 的出现次数
	windowOneCount := 0
	// 记录结果长度
	res := 0

	// 开始滑动窗口模板
	for right < len(nums) {
		// 扩大窗口
		if nums[right] == 1 {
			windowOneCount++
		}
		right++

		for right-left-windowOneCount > k {
			// 当窗口中需要替换的 0 的数量大于 k，缩小窗口
			if nums[left] == 1 {
				windowOneCount--
			}
			left++
		}
		// 此时一定是一个合法的窗口，求最大窗口长度
		res = max(res, right-left)
	}
	return res
}
func longestOnes2(nums []int, k int) int {

	maxNum := 0
	useK := 0
	left, right := 0, 0

	for right < len(nums) {
		c := nums[right]
		if c == 0 {
			// 开始借用，没有控制多借用了
			useK++
		}
		right++

		// 此时借多了，就需要换回来
		for useK > k {
			d := nums[left]
			if d == 0 {
				useK--
			}
			left++
		}

		// 保证不会多借用
		maxNum = max(maxNum, right-left)
	}

	return maxNum
}

func longestOnes(nums []int, k int) int {

	maxNum := 0
	useK := 0
	left, right := 0, 0

	for right < len(nums) {
		c := nums[right]
		if c == 0 && useK < k {
			useK++
		}
		right++
		maxNum = max(maxNum, right-left)

		for useK == k {
			d := nums[left]
			if d == 0 {
				useK--
			}
			left++
		}
	}

	return maxNum
}
