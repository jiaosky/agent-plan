package main

import "fmt"

func main() {
	//piles := []int{3, 6, 7, 11}
	//h := 8

	piles := []int{30, 11, 23, 4, 20}
	h := 5
	fmt.Println(minEatingSpeed(piles, h))
}

// 计算传入的速度，吃完全部要多久
func f(piles []int, speed int) int {
	hour := 0
	for i := range piles {
		hour += piles[i] / speed
		if piles[i]%speed > 0 {
			hour++
		}
	}
	return hour
}

func minEatingSpeed(piles []int, h int) int {
	// 速度最小值是 1 ，最大值就是香蕉的数量
	// [left,right）
	left, right := 1, 10<<9+1

	for left < right {
		// 居中的速度
		mid := left + (right-left)/2
		// 小于目标时间，说明速度太快，需要降低速度
		speed := f(piles, mid)
		if speed == h {
			// 等于目标时间，但是不一定是最小，继续找
			right = mid
		} else if speed < h {
			right = mid
		} else if speed > h {
			// 大于目标时间，说明速度太低，需要降低速度
			left = mid + 1
		}
	}

	return left
}
