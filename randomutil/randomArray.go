package randomutil

import (
	"math/rand"
	"time"
)

// 随机数种子源
var rd = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomIntArray 生成指定长度的int类型的随机数组
// @param minValue 最小值
// @param maxValue 最大值
// @param length 数组长度
// @return []int 生成的随机数组
func RandomIntArray(minValue int, maxValue int, length int) []int {
	if length <= 0 {
		return []int{}
	}
	// 计算随机数的有效范围
	rangeSize := maxValue - minValue
	// 创建随机数组
	var randomArray = make([]int, length)
	for index := range randomArray {
		randomArray[index] = minValue + rd.Intn(rangeSize)
	}
	return randomArray
}
