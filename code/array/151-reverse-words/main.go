package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(reverseWords("the sky is blue"))
	fmt.Println(reverseWords("  hello world  "))
	fmt.Println(reverseWords("a good   example"))
}

func reverseWords(s string) string {

	splits := splitWords(s)

	for i := 0; i < len(splits)/2; i++ {
		splits[i], splits[len(splits)-i-1] = splits[len(splits)-i-1], splits[i]
	}

	return strings.Join(splits, " ")
}

func splitWords(s string) []string {
	var words []string
	t := ""
	for i := 0; i < len(s); i++ {
		if s[i:i+1] != " " && s[i:i+1] != "" {
			t += s[i : i+1]
		} else {
			if len(t) > 0 {
				words = append(words, t)
			}
			t = ""
		}
	}
	if len(t) > 0 {
		words = append(words, t)
	}

	return words
}
