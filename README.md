# mineor.xyz


hugo 文章github备份

网站地址：https://www.mineor.xyz/


package matrix

import "fmt"

func printMatrix(m [][]int) {
	for _, row := range m {
		for _, col := range row {
			fmt.Printf("%d,", col)
		}
		fmt.Printf("\n")
	}
}
func InClockRotate(m [][]int) [][]int {
	d := len(m)
	for _, row := range m {
		for i := 0; i < d/2; i++ {
			tmp := row[i]
			row[i] = row[d-i-1]
			row[d-i-1] = tmp
		}
	}
	for row := 0; row < d; row++ {
		for col := 0; col < d-row-1; col++ {
			m[row][col], m[d-col-1][d-row-1] = m[d-col-1][d-row-1], m[row][col]
		}
	}
	printMatrix(m)
	return m
}

func Entry() {
	var m [][]int
	m = append(m, []int{1, 2}, []int{3, 4})
	//m = append(m, []int{1, 2, 3}, []int{4, 5, 6}, []int{7, 8, 9})
	//m = append(m, []int{1, 2, 3, 4}, []int{5, 6, 7, 8}, []int{9, 10, 11, 12}, []int{13, 14, 15, 16})
	InClockRotate(m)
}


package tree

type Tree struct {
	Val   int
	Left  *Tree
	Right *Tree
}


func BuildOrderTree(data []int) *Tree {
	var tree *Tree
	for _, d := range data {
		tree = insertNodeToTreeNoRec(tree, d)
	}

	return tree
}

func insetNodeToTree(t *Tree, val int) *Tree {
	if t == nil {
		return &Tree{Val: val}
	}
	if val < t.Val {
		t.Left = insetNodeToTree(t.Left, val)
	} else {
		t.Right = insetNodeToTree(t.Right, val)
	}
	return t
}

func insertNodeToTreeNoRec(t *Tree, val int) *Tree {
	newNode := &Tree{Val: val}
	if t == nil {
		return newNode
	}
	cur := t
	for cur != nil {
		if val < cur.Val {
			if cur.Left == nil {
				cur.Left = newNode
				break
			} else {
				cur = cur.Left
			}
		} else {
			if cur.Right == nil {
				cur.Right = newNode
				break
			} else {
				cur = cur.Right
			}
		}
	}
	return t
}

package tree

import "fmt"

func PrintInOrder(t *Tree) {
	if t != nil {
		PrintInOrder(t.Left)
		fmt.Println(t.Val)
		PrintInOrder(t.Right)
	}
}

func PrintPreOrder(t *Tree) {
	if t != nil {
		fmt.Println(t.Val)
		PrintPreOrder(t.Left)
		PrintPreOrder(t.Right)
	}
}

func PrintPostOrder(t *Tree) {
	if t != nil {
		PrintPostOrder(t.Left)
		PrintPostOrder(t.Right)
		fmt.Println(t.Val)
	}
}
func PrintInLevel(t *Tree) [][]int {
	var level []*Tree
	var levelResult [][]int
	var currentLevelResult []int
	if t == nil {
		return [][]int{}
	}
	cur := t
	level = append(level, cur) // 根节点入
	for len(level) > 0 {
		currentLevelLength := len(level)
		for i := 0; i < currentLevelLength; i++ {
			popped := level[i]
			currentLevelResult = append(currentLevelResult, popped.Val)
			if popped.Left != nil {
				level = append(level, popped.Left)
			}
			if popped.Right != nil {
				level = append(level, popped.Right)
			}
		}
		level = level[currentLevelLength:]
		levelResult = append(levelResult, currentLevelResult)
		currentLevelResult = []int{}
	}
	return levelResult
}

