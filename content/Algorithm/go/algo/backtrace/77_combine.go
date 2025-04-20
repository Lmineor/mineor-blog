package backtrace

//给定两个整数 n 和 k，返回 1 ... n 中所有可能的 k 个数的组合。
//
//示例: 输入: n = 4, k = 2 输出: [ [2,4], [3,4], [2,3], [1,2], [1,3], [1,4], ]
//
//#

var result [][]int
var path []int

func Combine(n int, k int) [][]int {
	result = make([][]int, 0)
	path = make([]int, 0)
	backTrace(n, k, 1)
	return result
}

func backTrace(n int, k int, start int) {
	if len(path) == k {
		tmp := make([]int, k)
		copy(tmp, path)
		result = append(result, tmp)
		return
	}
	for i := start; i <= n; i++ {
		path = append(path, i)
		backTrace(n, k, i+1)
		path = path[:len(path)-1]
	}
}
