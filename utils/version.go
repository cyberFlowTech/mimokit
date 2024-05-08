package utils

import (
	"strconv"
	"strings"
)

func VersionCompare(v1, v2 string) int {
	// 将版本号分割为数组
	v1Arr := strings.Split(v1, ".")
	v2Arr := strings.Split(v2, ".")

	// 遍历数组，逐位比较
	for i := 0; i < len(v1Arr) && i < len(v2Arr); i++ {
		v1Num, err := strconv.Atoi(v1Arr[i])
		if err != nil {
			return -1
		}
		v2Num, err := strconv.Atoi(v2Arr[i])
		if err != nil {
			return -1
		}

		if v1Num != v2Num {
			return v1Num - v2Num
		}
	}

	// 如果数组长度不等，则比较数组长度
	if len(v1Arr) < len(v2Arr) {
		return -1
	} else if len(v1Arr) > len(v2Arr) {
		return 1
	} else {
		return 0
	}
}
