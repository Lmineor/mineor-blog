package main

import "fmt"

type TreeNode struct {
	Val int
	Left *TreeNode
	Right *TreeNode
}
func PreOrder(root *TreeNode) []int {
	result := make([]int, 0)
	var f func(t *TreeNode)
	f = func(t *TreeNode) {
		if t != nil {
			result = append(result, t.Val)
			f(t.Left)
			f(t.Right)
		}
	}
	f(root)
	return result
}

func main(){
	t := &TreeNode{
		Val:   1,
		Left:  &TreeNode{
			Val:   2,
			Left:  &TreeNode{Val: 4},
			Right: &TreeNode{Val: 5},
		},
		Right: &TreeNode{
			Val:   3,
			Left:  &TreeNode{Val: 6},
			Right: &TreeNode{Val: 7},
		},
	}
	fmt.Println(PreOrder(t))
}