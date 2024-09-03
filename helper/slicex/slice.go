package slicex

import (
	"errors"
	"fmt"
	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/proto"
	"reflect"
	"runtime"
	"time"
)

// 过滤切片中所有值等于k的元素并返回新切片
func FilterSlice[T comparable](s []T, k T) []T {
	var res []T
	for _, v := range s {
		if v != k {
			res = append(res, v)
		}
	}
	return res
}

func Contains[T comparable](s []T, k T) bool {
	for _, v := range s {
		if v == k {
			return true
		}
	}
	return false
}

func BSearch[T any](arr []T, cmp func(ele T) int) int {
	i, j := 0, len(arr)
	for i < j {
		h := int(uint(i+j) >> 1)
		x := cmp(arr[h])
		if x == 0 {
			return h
		}
		if x < 0 { // < 0左边
			j = h
		} else { // > 0右边
			i = h + 1
		}
	}
	return -1
}

func Between[T constraints.Ordered](a, lhs, rhs T) bool {
	if lhs > rhs {
		lhs, rhs = rhs, lhs
	}
	return a >= lhs && a <= rhs
}

// val in [pivot-fact, pivot+fact]
func LooseEqual[T constraints.Integer](val, pivot, fact T) bool {
	if pivot-fact > pivot || pivot+fact < pivot {
		return false
	}
	return pivot-fact <= val && val <= pivot+fact
}

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
	if len(s.items) == 0 {
		panic("Stack is empty")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func ForeachFn(fn ...func() error) error {
	for i, v := range fn {
		if err := v(); err != nil {
			pc := reflect.ValueOf(v).Pointer()
			funcName := runtime.FuncForPC(pc).Name()
			return fmt.Errorf("%d func: %s err: %v", i, funcName, err)
		}
	}
	return nil
}

var Continue = fmt.Errorf("continue")

func Retry(n int) func(f func() error) error {
	if n == 0 {
		panic("n must be greater than 0")
	}
	return func(f func() error) error {
		el := make([]error, 0)
		for i := 0; i < n; i++ {
			if err := f(); err == nil {
				return nil
			} else if errors.Is(err, Continue) {
				continue
			} else {
				el = append(el, err)
			}
		}
		if len(el) == 0 {
			return nil
		}
		return fmt.Errorf("%+v", el)
	}
}

func RetrySleep(n int, d time.Duration) func(f func() error) error {
	if n == 0 {
		panic("n must be greater than 0")
	}
	return func(f func() error) error {
		el := make([]error, 0)
		for i := 0; i < n; i++ {
			if err := f(); err == nil {
				return nil
			} else if errors.Is(err, Continue) {
				continue
			} else {
				time.Sleep(d)
				el = append(el, err)
			}
		}
		if len(el) == 0 {
			return nil
		}
		return fmt.Errorf("%+v", el)
	}
}

/*
func ConvTo[T, F any](list []F) []T {
	result := make([]T, len(list))
	for i := range list {
		result[i] = any(list[i]).(T)
	}
	return result
}

func Transform[T, U any](list []T, f func(T) U) []U {
	result := make([]U, len(list))
	for i := range list {
		result[i] = f(list[i])
	}
	return result
}
*/

func RemoveUint32[T comparable](s []T, v T) []T {
	for i, val := range s {
		if val == v {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func DeepCopyProtoSlice[T proto.Message](list []T) []T {
	result := make([]T, len(list))
	for i := range list {
		result[i] = proto.Clone(list[i]).(T)
	}
	return result
}
