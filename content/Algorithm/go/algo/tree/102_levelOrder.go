package tree

//给你一个二叉树，请你返回其按 层序遍历 得到的节点值。 （即逐层地，从左到右访问所有节点）。

func levelOrder(root *TreeNode) [][]int {
	var result [][]int
	if root == nil {
		return result
	}
	var level []*TreeNode
	level = append(level, root) // 根节点入队
	var currentLevelResult []int
	var currentLevelLen int
	for len(level) > 0 {
		currentLevelResult = []int{}
		currentLevelLen = len(level)
		for i := 0; i < currentLevelLen; i++ {
			node := level[0]
			currentLevelResult = append(currentLevelResult, node.Val)
			level = level[1:]
			if node.Left != nil {
				level = append(level, node.Left)
			}
			if node.Right != nil {
				level = append(level, node.Right)
			}
		}
		result = append(result, currentLevelResult)
	}
	return result
}
