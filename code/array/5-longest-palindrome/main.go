package main

import "fmt"

func main() {

	fmt.Println(longestPalindrome("babad"))
	fmt.Println(longestPalindrome("cbbd"))

}

func longestPalindrome(s string) string {

	var ss string

	for i := range s {
		s1 := sub(s, i, i)
		s2 := sub(s, i, i+1)

		m := s1
		if len(m) < len(s2) {
			m = s2
		}

		if len(m) > len(ss) {
			ss = m
		}
	}

	return ss
}

func sub(s string, centerLeft, centerRight int) string {

	for centerLeft >= 0 && centerRight < len(s) && s[centerLeft] == s[centerRight] {
		centerLeft--
		centerRight++
	}
	return s[centerLeft+1 : centerRight]

}
