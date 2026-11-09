package main

import "fmt"

func main() {
	ls := build([]int{1, 3, 3, 5, 5, 7, 8, 9})

	duplicates := deleteDuplicates(ls)

	printNode(duplicates)

}

func printNode(nodes *ListNode) {
	for nodes != nil {
		fmt.Print(nodes.Val, " ")
		nodes = nodes.Next
	}

}

func build(nums []int) *ListNode {

	if nums == nil || len(nums) == 0 {
		return nil
	}

	head := &ListNode{
		Val: nums[0],
	}
	cur := head

	for i := 1; i < len(nums); i++ {
		cur.Next = &ListNode{Val: nums[i]}
		cur = cur.Next
	}

	return head
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func deleteDuplicates(head *ListNode) *ListNode {

	if head == nil {
		return nil
	}
	slow, fast := head, head

	for fast != nil {
		if fast.Val != slow.Val {
			slow.Next = fast
			slow = slow.Next
		}
		fast = fast.Next
	}
	slow.Next = nil
	return head
}
