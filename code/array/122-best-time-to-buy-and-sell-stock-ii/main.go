package main

import (
	"fmt"
	"math"
)

func main() {
	price := []int{2, 1, 2, 0, 1}
	k := maxProfit2(price)
	fmt.Println(k)
}

func maxProfit2(prices []int) int {
	minPrice := math.MaxInt
	maxprofit := 0
	sumMaxProfit := 0
	for i := 0; i < len(prices); i++ {
		if maxprofit == 0 && prices[i] < minPrice { //  利润为空代表新的一轮，找最小的
			minPrice = prices[i]
		} else {
			// 历史最大利润 > 当天的利润
			// 说明昨天应该卖出，今天买入
			// 然后利润变成 0
			if maxprofit > prices[i]-minPrice {
				sumMaxProfit = sumMaxProfit + maxprofit
				maxprofit = 0
				minPrice = prices[i]
				continue
			} else {
				maxprofit = prices[i] - minPrice
			}
		}
	}

	sumMaxProfit = sumMaxProfit + maxprofit

	return sumMaxProfit
}
