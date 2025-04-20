package backtrace

import "fmt"

//找出所有相加之和为 n 的 k 个数的组合，且满足下列条件：
//
//只使用数字1到9
//每个数字 最多使用一次
//返回 所有可能的有效组合的列表 。该列表不能包含相同的组合两次，组合可以以任何顺序返回。

var result3 [][]int
var path3 []int

func combinationSum3(k int, n int) [][]int {
	result3, path3 = make([][]int, 0), make([]int, 0, k)
	bt1(k, n, 1)
	return result3
}

func bt1(k int, n int, start int) {
	if len(path3) == k && sum(path3) == n {
		tmp := make([]int, k)
		copy(tmp, path3)
		result3 = append(result3, tmp)
		return
	}
	s := 0
	if n <= 9{
		 s = n
	}else{
		s=9
	}
	for i := start; i <= s; i++ {
		path3 = append(path3, i)
		bt1(k, n, i+1)
		path3 = path3[:len(path3)-1]
	}
}

func sum(l []int) int {
	r := 0
	for _, v := range l {
		r += v
	}
	return r
}

func main() {
	fmt.Println(combinationSum3(4, 1))
}
