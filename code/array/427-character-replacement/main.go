package main

import "fmt"

func main() {
	//fmt.Println(characterReplacement("ABAB", 2))
	fmt.Println(characterReplacement2("AABABBA", 1))
	//fmt.Println(characterReplacement("ABBB", 2))
}

func characterReplacement2(s string, k int) int {
	left, right := 0, 0
	maxCount := 0
	maxLength := 0

	cnt := [26]int{}

	for right < len(s) {
		// 记录窗口中元素出现的次数
		cnt[s[right]-'A']++
		// 找最大数量
		maxCount = max(cnt[s[right]-'A'], maxCount)
		right++

		// 窗口元素的数量 - 最大的数量 > 可以替换的数量
		for right-left-maxCount > k {
			cnt[s[left]-'A']--
			left++
		}
		maxLength = max(maxLength, right-left)
	}

	return maxLength
}

func characterReplacement(s string, k int) int {
	left, right := 0, 0
	maxCount := 0
	maxLength := 0

	m := make(map[byte]int)

	for right < len(s) {
		// 记录窗口中元素出现的次数
		m[s[right]]++
		// 找最大数量
		maxCount = max(m[s[right]], maxCount)
		right++

		// 窗口元素的数量 - 最大的数量 > 可以替换的数量
		for right-left-maxCount > k {
			m[s[left]]--
			left++
		}
		maxLength = max(maxLength, right-left)
	}

	return maxLength
}
