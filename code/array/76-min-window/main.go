package main

import (
	"fmt"
	"math"
)

func main() {
	s := "aa"
	t := "aa"
	fmt.Println(minWindow(s, t))
}

func minWindow(s string, t string) string {

	left, right := 0, 0
	start := 0
	length := math.MaxInt

	need := make(map[byte]int)
	find := make(map[byte]int)
	matchNum := 0

	for _, c := range []byte(t) {
		need[c]++
	}

	for right < len(s) {
		c := s[right]
		right++

		if _, ok := need[c]; ok {
			find[c]++ // 放入数量不限制
			if find[c] == need[c] {
				matchNum++ // 放入相同数量时才标记
			}
		}

		// 命中的元素 == 需要的元素
		for matchNum == len(need) {
			if right-left < length {
				start = left
				length = right - start
			}

			d := s[left]
			left++

			if _, ok := need[d]; ok {
				if find[d] == need[d] {
					matchNum-- // 移除前数量相同，移除后就不匹配了
				}
				find[d]-- // 移除数量不限制
			}
		}
	}

	if length == math.MaxInt {
		return ""
	}

	return s[start : start+length]
}

func minWindow3(s string, t string) string {
	need, window := make(map[byte]int), make(map[byte]int)
	for i := range t {
		need[t[i]]++
	}

	left, right := 0, 0
	valid := 0
	// 记录最小覆盖子串的起始索引及长度
	start, length := 0, math.MaxInt32
	for right < len(s) {
		// c 是将移入窗口的字符
		c := s[right]
		// 扩大窗口
		right++
		// 进行窗口内数据的一系列更新
		if _, ok := need[c]; ok {
			window[c]++
			if window[c] == need[c] {
				valid++
			}
		}

		// 判断左侧窗口是否要收缩
		for valid == len(need) {
			// 在这里更新最小覆盖子串
			if right-left < length {
				start = left
				length = right - left
			}
			// d 是将移出窗口的字符
			d := s[left]
			// 缩小窗口
			left++
			// 进行窗口内数据的一系列更新
			if _, ok := need[d]; ok {
				if window[d] == need[d] {
					valid--
				}
				window[d]--
			}
		}
	}
	// 返回最小覆盖子串
	if length == math.MaxInt32 {
		return ""
	}
	return s[start : start+length]
}
func minWindow2(s string, t string) string {

	needWindows := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		needWindows[t[i]] = needWindows[t[i]] + 1
	}

	findWindows := make(map[byte]int)

	left, right := 0, 0
	n := len(s)

	length := math.MaxInt
	minString := s

	valid := 0

	for right < n {

		// 放入的元素
		c := s[right]
		right++

		// 是需要的字符
		if _, ok := needWindows[c]; ok {
			findWindows[c] = findWindows[c] + 1
			if findWindows[c] == needWindows[c] {
				valid++ // 找到了一个
			}
		}

		// 全部找到
		for valid == len(needWindows) {
			// 这里就是本轮结束
			if length > len(s[left:right]) {
				minString = s[left:right]
				length = len(minString)
			}

			d := s[left]
			left++
			if _, ok := needWindows[d]; ok {
				if findWindows[d] == needWindows[d] {
					// 不满足的数量+1
					valid--
				}
				findWindows[d] = findWindows[d] - 1
			}
		}

	}

	if length == math.MaxInt {
		return ""
	}

	return minString
}
