package helper

import (
	"context"
	"github.com/shopspring/decimal"
	"math"
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

func DivisionFloat64(dividend, divisor float64) float64 {
	if divisor == 0 {
		return 0
	} else {
		return dividend / divisor
	}
}

// float64 保留3位小数
func Float64DecimalRound3(value float64) float64 {
	res, _ := decimal.NewFromFloat(value).Round(5).Float64()
	if res > 0.001 {
		res, _ := decimal.NewFromFloat(value).Round(3).Float64()
		return res
	} else if res > 0 {
		return 0.001
	} else {
		return 0
	}
}

// float64 除法，保留3位小数
func DivisionFloat64DecimalRound3(dividend, divisor float64) float64 {
	if divisor == 0 {
		return 0
	} else {
		return Float64DecimalRound3(dividend / divisor)
	}
}

// float64 除法，保留n位小数(少数需要精确度高的地方使用,建议少用)
func DivisionFloat64DecimalRoundN(dividend, divisor float64, round int32) float64 {
	if round > 10 {
		round = 10
	}

	if divisor == 0 {
		return 0
	} else {
		value := dividend / divisor
		base := 1 / float64(math.Pow10(int(round)))
		res, _ := decimal.NewFromFloat(value).Round(round + 2).Float64()
		if res > base {
			res, _ := decimal.NewFromFloat(value).Round(round).Float64()
			return res
		} else if res > 0 {
			return base
		} else {
			return 0
		}
	}
}

// 使用泛型检查 slice 中是否包含某个元素
func IsExistInList[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// 返回字符串的字符长度
func GetCharacterCount(str string) int {
	return utf8.RuneCountInString(str)
}
