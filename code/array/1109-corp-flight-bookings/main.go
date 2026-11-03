package main

import "fmt"

func main() {
	bookings := [][]int{{1, 2, 10}, {2, 2, 15}}
	n := 2
	fmt.Print(corpFlightBookings(bookings, n))
}

// 这里有 n 个航班，它们分别从 1 到 n 进行编号。
//
// 有一份航班预订表 bookings ，表中第 i 条预订记录 bookings[i] = [firsti, lasti, seatsi] 意味着在从 firsti 到 lasti （包含 firsti 和 lasti ）的 每个航班 上预订了 seatsi 个座位。
//
// 请你返回一个长度为 n 的数组 answer，里面的元素是每个航班预定的座位总数。
func corpFlightBookings(bookings [][]int, n int) []int {

	answer := make([]int, n)
	diff := make([]int, n)
	for _, booking := range bookings {
		first := booking[0] - 1
		last := booking[1] - 1
		seat := booking[2]

		diff[first] = diff[first] + seat
		if last+1 < n {
			diff[last+1] = diff[last+1] - seat
		}
	}
	fmt.Println(diff)
	answer[0] = diff[0]
	for i := 1; i < len(diff); i++ {
		answer[i] = diff[i] + answer[i-1]
	}

	return answer
}
