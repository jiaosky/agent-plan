package main

import "fmt"

func main() {
	s := "abab"
	p := "ab"
	fmt.Println(findAnagrams(s, p))
}

func findAnagrams(s, p string) (ans []int) {
	// 统计 p 的每种字母的出现次数
	cntP := [26]int{}
	for _, c := range p {
		cntP[c-'a']++
	}

	cntS := [26]int{} // 统计 s 的长为 len(p) 的子串 t 的每种字母的出现次数
	for right, c := range s {
		cntS[c-'a']++ // 右端点字母进入窗口
		left := right - len(p) + 1
		if left < 0 { // 窗口长度不足 len(p)
			continue
		}
		if cntS == cntP { // t 和 p 的每种字母的出现次数都相同
			ans = append(ans, left) // t 左端点下标加入答案
		}
		cntS[s[left]-'a']-- // 左端点字母离开窗口
	}
	return
}

func findAnagrams3(s string, p string) []int {

	need := make(map[byte]int)
	find := make(map[byte]int)

	left, right, matchNum := 0, 0, 0
	res := make([]int, 0)

	for _, c := range p {
		need[byte(c)]++
	}

	for right < len(s) {
		c := s[right]
		right++
		if _, ok := need[c]; ok {
			find[c]++
			if find[c] == need[c] {
				matchNum++
			}
		}

		// 要找到全部的子串，所以要遍历完一遍呢
		for right-left >= len(p) {
			if matchNum == len(need) {
				res = append(res, left)
			}

			d := s[left]
			left++
			if _, ok := need[d]; ok {
				if find[d] == need[d] {
					matchNum--
				}
				find[d]--
			}
		}
	}
	return res
}

func findAnagrams2(s string, p string) []int {

	need := make(map[byte]int)
	find := make(map[byte]int)

	left, right, matchNum := 0, 0, 0
	res := make([]int, 0)

	for _, c := range p {
		need[byte(c)]++
	}

	for right < len(s) {
		c := s[right]
		right++
		if _, ok := need[c]; ok {
			find[c]++
			if find[c] == need[c] {
				matchNum++
			}
		}

		// 要找到全部的子串，所以要遍历完一遍呢
		for matchNum == len(need) {
			if len(p) == right-left {
				res = append(res, left)
			}

			d := s[left]
			left++
			if _, ok := need[d]; ok {
				if find[d] == need[d] {
					matchNum--
				}
				find[d]--
			}
		}
	}
	return res
}
