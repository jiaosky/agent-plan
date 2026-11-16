package main

import (
	"fmt"
)

func main() {
	s := ".,"
	fmt.Println(isPalindrome(s))
}

func isalnum(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9')
}

func isPalindrome(s string) bool {
	left, right := 0, len(s)-1
	for left < right {
		for !isalnum(s[left]) && left < right {
			left++
		}

		for !isalnum(s[right]) && right > left {
			right--
		}

		if left < right {
			if s[left] != s[right] {
				return false
			}
			left++
			right--
		}
	}

	return true
}
