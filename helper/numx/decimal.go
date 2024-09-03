package numx

import (
	"github.com/shopspring/decimal"
	"math"
)

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
