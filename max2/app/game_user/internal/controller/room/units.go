package room

import "github.com/gogf/gf/v2/util/grand"

// Probability 判断指定百分比概率是否触发
// percent: 0-100 的百分比
func Probability(percent int) bool {
	if percent <= 0 {
		return false
	}
	if percent >= 100 {
		return true
	}
	return grand.Intn(100) < percent
}
