package helper

import (
	"context"
	"strconv"
	"time"
	"unicode/utf8"
)

func GetCtxTimeOut(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func GetCtxTimeOut10() (context.Context, context.CancelFunc) {
	return GetCtxTimeOut(10 * time.Second)
}

func IntToString(value int) string {
	return strconv.Itoa(value)
}

func Uint64ToString(value uint64) string {
	return strconv.FormatUint(value, 10)
}

func StringToInt(value string) int {
	data, _ := strconv.Atoi(value)
	return data
}

func String2Int(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func StringToUint64(value string) uint64 {
	data, _ := strconv.ParseUint(value, 10, 64)
	return data
}

func Float64ToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func Float32ToString(value float32) string {
	return strconv.FormatFloat(float64(value), 'f', -1, 32)
}

// 返回字符串的字符长度
func GetCharacterCount(str string) int {
	return utf8.RuneCountInString(str)
}
