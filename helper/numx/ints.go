package numx

import (
	. "golang.org/x/exp/constraints"
	"strconv"
)

func I2A[T Signed](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

func U2A[T Unsigned](i T) string {
	return strconv.FormatUint(uint64(i), 10)
}

func A2I[T Signed](s string) (T, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return T(i), nil
}

func A2IOr[T Signed](s string, or T) T {
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return T(i)
	}
	return or
}

func A2UOr[T Unsigned](s string, or T) T {
	i, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		return T(i)
	}
	return or
}

func A2U[T Unsigned](s string) (T, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return T(i), nil
}

func IConvU32(i int32) uint32 {
	if i >= 0 {
		return uint32(i)
	}
	return uint32(-i)
}

func Or[T Integer](a, b T) T {
	if a != 0 {
		return a
	}
	return b
}

func Clamp[T Integer](v, min, max T) T {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func OrElse[T Integer](v *T, vv T) {
	if v != nil && *v == 0 {
		*v = vv
	}
}

type Pair[T, U any] struct {
	First  T
	Second U
}

func NewPair[T, U any](f T, s U) *Pair[T, U] {
	return &Pair[T, U]{First: f, Second: s}
}

func ToPair[K comparable, V any](m map[K]V) []*Pair[K, V] {
	list := make([]*Pair[K, V], 0, len(m))
	for k, v := range m {
		list = append(list, &Pair[K, V]{First: k, Second: v})
	}
	return list
}
