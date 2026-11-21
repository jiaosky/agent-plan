package main

import "fmt"

func main() {
	//fmt.Println(lengthOfLongestSubstring("abcabcbb"))
	//fmt.Println(lengthOfLongestSubstring("bbbbb"))
	fmt.Println(lengthOfLongestSubstring("pwwkew"))
}

func lengthOfLongestSubstring(s string) int {

	find := make(map[string]int)
	left, right := 0, 0
	maxLen := 0
	for right < len(s) {
		c := s[right : right+1]
		right++
		find[c]++

		for find[c] > 1 {
			d := s[left : left+1]
			left++
			find[d]--
		}
		if maxLen < right-left {
			maxLen = right - left
		}
	}

	return maxLen

}

func lengthOfLongestSubstring2(s string) int {
	find := make(map[byte]int)

	left, right := 0, 0
	subLength := 0
	maxSubLength := 0

	for right < len(s) {
		c := s[right]
		right++

		for find[c] < 1 {
			find[c]++
			subLength++
		}

		// right 代表的字符遇到了重复的
		// 子串长度 < 窗口数，才说明有重复的
		for subLength < (right-left) && left < right {
			d := s[left]

			// 前后相同，则长度不变，继续就行
			if c == d {
				left++
			} else {
				// 前后不同，说明是子串中间有重复的，此时 subLength 就少了，类似 abcb ,rigth 是 b，left 是 a
				find[d]--
				subLength--
				left++ // 继续缩小窗口
				// abb
			}
		}

		if subLength > maxSubLength {
			maxSubLength = subLength
		}
	}

	return maxSubLength
}
