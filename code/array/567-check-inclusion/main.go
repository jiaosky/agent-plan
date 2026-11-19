package main

import "fmt"

func main() {

	s1 := "ab"
	s2 := "eidbaooo"

	fmt.Println(checkInclusion(s1, s2))
}

func checkInclusion(s1 string, s2 string) bool {

	// 所谓排列，就是长度相同，字母出现的数量也相同

	need := make(map[byte]int)
	find := make(map[byte]int)
	for _, c := range s1 {
		need[byte(c)]++
	}

	left, right := 0, 0
	n := len(s2)

	matchNum := 0

	for right < n {
		c := s2[right]
		right++
		find[c]++

		if _, exist := need[c]; exist {
			if need[c] == find[c] {
				matchNum++
			}
		}

		// left ~ right 有足够的相同的字符了，现在看看是否满足长度相同
		for matchNum == len(need) {
			if right-left == len(s1) {
				return true
			}

			d := s2[left]
			left++
			if _, exist := need[d]; exist {
				if need[d] == find[d] {
					matchNum--
				}
				find[d]--
			}
		}
	}

	return false
}
