package utils

import "strings"

const chars = "L5MPd9wx2eWXfabc1gQijklmEF4UnoNOp3vyzABRh0CGHrs7tuIJKS6TDVq8YZ"

func ToBase62(num int64) string {
	result := ""
	for num > 0 {
		result = string(chars[num%62]) + result
		num /= 62
	}
	return result
}

func ToBase10(str string) int64 {
	var res int64 = 0
	for _, s := range str {
		index := strings.IndexRune(chars, s)
		res = res*62 + int64(index)
	}
	return res
}
