package main

import (
	"fmt"
	"math"
)

func main() {
	price := []int{2, 4, 11, 1, 100}
	k := maxProfit(price)
	fmt.Println(k)
}

// 121. 买卖股票的最佳时机
// 给定一个数组 prices ，它的第 i 个元素 prices[i] 表示一支给定股票第 i 天的价格。
// 你只能选择 某一天 买入这只股票，并选择在 未来的某一个不同的日子 卖出该股票。设计一个算法来计算你所能获取的最大利润。
// 返回你可以从这笔交易中获取的最大利润。如果你不能获取任何利润，返回 0 。
func maxProfit(prices []int) int {

	minprice := math.MaxInt // 最小值
	maxprofit := 0          // 最大差价

	for i := 0; i < len(prices); i++ {
		if prices[i] < minprice { // 当前价格最低
			minprice = prices[i]
		} else {
			// 看盈利状态
			if prices[i]-minprice > maxprofit {
				maxprofit = prices[i] - minprice
			}
		}
	}

	return maxprofit
}
