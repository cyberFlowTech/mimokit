package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var randomStringSeed int64

// RandomString 返回一个 chars 组成的长度为 n 的随机字符串
func RandomString(n int, chars string) string {
	if randomStringSeed == 0 {
		randomStringSeed = time.Now().UnixNano()
		rand.Seed(randomStringSeed)
	}
	var letters = []rune(chars)
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// RandomAlphanumeric 返回一个由大小写字母和数字组成的长度为 n 的随机字符串
func RandomAlphanumeric(n int) string {
	return RandomString(n, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}

// RandomCode 返回一个高辨识度的随机字符串，去除了容易引起视觉混淆的字母（如 L, I, 1, O, 0）。可以用于邀请码等场景。
func RandomCode() string {
	return RandomNCode(7)
}

// RandomNCode 返回一个高辨识度的长度为N的随机字符串，去除了容易引起视觉混淆的字母（如 L, I, 1, O, 0）。
func RandomNCode(n int) string {
	return RandomString(n, "ABCDEFGHKMNPQRSTWXYZ23456789")
}

func RandomFileCode() string {
	return fmt.Sprintf("%d%s", UnixMilliNow(), RandomNCode(7))
}

func EncodeCode(n uint64) string {
	strs := "ABCDEFGHKMNPQRSTWXYZ23456789"
	result := []byte{}
	for n > 0 {
		result = append(result, strs[n%uint64(len(strs))])
		n = n / uint64(len(strs))
	}
	return string(result)
}

func EncodeSnToCode(sn string) string {
	if s, err := strconv.ParseUint(sn, 10, 64); err != nil {
		return sn
	} else {
		return EncodeCode(s)
	}
}
