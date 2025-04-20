package backtrace

// 给你一个字符串 s，请你将 s 分割成一些 子串，使每个子串都是 回文串 。返回 s 所有可能的分割方案。
//
// 示例 1：
//
// 输入：s = "aab"
// 输出：[["a","a","b"],["aa","b"]]
// 示例 2：
//
// 输入：s = "a"
// 输出：[["a"]]
var result1 [][]string
var path1 []string

func partition(s string) [][]string {
	result1 = make([][]string, 0)
	path1 = make([]string, 0)
	backTractPartition(s, 0)
	return result1
}

func backTractPartition(s string, start int) {
	if start == len(s) { // 起始位置等于s的大小，说明已经找到一组方案
		tmp := make([]string, len(path1))
		copy(tmp, path1)
		result1 = append(result1, tmp)
		return
	}
	for i := start; i < len(s); i++ {
		str := s[start : i+1]
		if isHuiwen(str) {
			// 如果是回文字符串，接下来找后面符合的字符串
			path1 = append(path1, str)
			backTractPartition(s, i+1)
			path1 = path1[:len(path1)-1]
		}
	}
}

func isHuiwen(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-i-1] {
			return false
		}
	}
	return true
}
