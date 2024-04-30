package randomutil

import (
	"math/rand"
	"time"
	"unsafe"
)

const (
	// 定义可用字符的字符串
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	// 62个字符能用6个比特位数表示完
	letterIdBits = 6
	// 最大的字母id掩码
	letterIdMask = 1<<letterIdBits - 1
	// 可用次数的最大值
	letterIdMax = letterIdMask / letterIdBits
)

// 随机数种子源
var random = rand.NewSource(time.Now().UnixNano())

// RandomString 生成指定长度的随机字符串
// @param length 字符串长度
// @return string 生成的随机字符串
func RandomString(length int) string {
	if length <= 0 {
		return ""
	}
	// 创建一个长度为 length 的字节切片
	bytes := make([]byte, length)
	// 循环生成随机字符串
	for i, cache, remain := length-1, random.Int63(), letterIdMax; i >= 0; {
		// 检查随机数生成器是否用尽所有随机数
		if remain == 0 {
			cache, remain = random.Int63(), letterIdMax
		}
		// 从可用字符的字符串中随机选择一个字符
		if idx := int(cache & letterIdMask); idx < len(letters) {
			bytes[i] = letters[idx]
			i--
		}
		// 右移比特位数，为下次选择字符做准备
		cache >>= letterIdBits
		remain--
	}
	// 将字节切片转换为字符串并返回
	return *(*string)(unsafe.Pointer(&bytes))
}
