package random

import (
	"math/rand"
	"time"
)

/*
	参　数：
		min: 最小值,max: 最大值
	返回值：
		int64: 生成的随机数
*/
func RandInt(min, max int) int {
	if min >= max || max == 0 {
		return max
	}
	time.Sleep(time.Millisecond * 3)
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
