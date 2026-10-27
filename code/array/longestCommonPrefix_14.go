package main

import "fmt"

/*
*
14. 最长公共前缀
编写一个函数来查找字符串数组中的最长公共前缀。
如果不存在公共前缀，返回空字符串 ""。

示例 1：

输入：strs = ["flower","flow","flight"]
输出："fl"
示例 2：

输入：strs = ["dog","racecar","car"]
输出：""
解释：输入不存在公共前缀。
*/

func main() {
	strs := []string{"flower", "flower", "flower", "flower"}
	rs := longestCommonPrefix(strs)
	fmt.Println(rs)
}

// 纵向
// 依次对比每列前面的内容是否相同
// 转化为二维数组
func longestCommonPrefix(strs []string) string {
	if len(strs) == 1 {
		return strs[0]
	}
	// y 轴
	for col := 0; col < len(strs[0]); col++ {
		// x 轴
		for index := 1; index < len(strs); index++ {
			if col == len(strs[index]) || strs[index][col] != strs[0][col] {
				// 中间匹配
				return strs[0][:col]
			}
		}
	}

	return strs[0]
}

func minLength(a string, b string) int {
	if len(a) > len(b) {
		return len(b)
	}
	return len(a)
}

// 横向
func longestCommonPrefixHengxiang(strs []string) string {
	if len(strs) == 1 {
		return strs[0]
	}
	str := strs[0]
	for i := 1; i < len(strs); i++ {
		next := strs[i]
		if len(next) == 0 {
			return ""
		}
		for j := minLength(str, next); j > 0; j-- {
			if str[:j] == next[:j] {
				str = next[:j]
				break
			}
			if j == 1 {
				return ""
			}
		}
	}

	return str
}

// 暴力
func longestCommonPrefixbaoli(strs []string) string {

	if len(strs) == 1 {
		return strs[0]
	}
	m := make(map[string]int)
	for _, v := range strs {
		for i := 0; i < len(v); i++ {
			if k, exist := m[v[0:i+1]]; exist {
				m[v[0:i+1]] = k + 1
			} else {
				m[v[0:i+1]] = 1
			}
		}
	}

	maxString := ""
	for k, v := range m {
		if v == len(strs) {
			if len(k) > len(maxString) {
				maxString = k
			}
		}
	}

	return maxString
}
