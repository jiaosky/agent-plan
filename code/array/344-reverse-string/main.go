package main

import "fmt"

func main() {
	b := []byte("hello")
	fmt.Println(b)
	reverseString(b)
	fmt.Println(b)

}

func reverseString(s []byte) {
	left, right := 0, len(s)-1
	for left <= right {
		s[left], s[right] = s[right], s[left]
		left++
		right--
	}
}
