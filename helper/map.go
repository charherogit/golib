package helper

import (
	"fmt"
	. "golang.org/x/exp/constraints"
)

func MergeMap[T comparable, U Integer](m ...map[T]U) map[T]U {
	if len(m) == 0 {
		return nil
	}
	mm := make(map[T]U)
	for _, v := range m {
		for kk, vv := range v {
			mm[kk] += vv
		}
	}
	return mm
}

func AddMap[T comparable, U Integer](m map[T]U, n ...map[T]U) {
	for _, v := range n {
		for kk, vv := range v {
			m[kk] += vv
		}
	}
}

func SameMapKey[T comparable, U any](m ...map[T]U) bool {
	if len(m) == 0 {
		return true
	}
	m0 := m[0]
	for _, v := range m[1:] {
		for k := range m0 {
			if _, ok := v[k]; !ok {
				return false
			}
		}
	}
	return true
}

func SliceToMap[T comparable, U any](slice []U, valueFn func(U) T) map[T]U {
	m := make(map[T]U, len(slice))
	for _, element := range slice {
		m[valueFn(element)] = element
	}
	return m
}

func EqualMap[T comparable, U Integer](m, n map[T]U) error {
	if len(m) != len(n) {
		return fmt.Errorf("req len(m): %d != ser len(n): %d map ser n: %v not equal req m: %v", len(m), len(n), n, m)
	}
	for k, v := range m {
		if n[k] != v {
			return fmt.Errorf("ser n[k]: %d != req m[k]: %d map ser n: %v not equal req m: %v", n[k], v, n, m)
		}
	}
	return nil
}

func ContainMap[T comparable, U Integer](m, n map[T]U) error {
	for k, v := range m {
		if n[k] < v {
			return fmt.Errorf("ser n[k]: %d < req m[k]: %d map ser n: %v not contain req m: %v", n[k], v, n, m)
		}
	}
	return nil
}

func GetMapKeys[T comparable, U any](m ...map[T]U) []T {
	var keys []T
	for _, v := range m {
		for k := range v {
			keys = append(keys, k)
		}
	}
	return keys
}

func SumMap[T comparable, U Integer](m map[T]U) U {
	var sum U
	for _, v := range m {
		sum += v
	}
	return sum
}

func DeleteMapKey[T comparable, U Integer](m, n map[T]U) {
	for k, v := range m {
		if n[k] >= v {
			delete(m, k)
		} else {
			m[k] -= n[k]
		}
	}
}

func DeleteMapWithSlice[T comparable, U comparable](m map[T]U, n []T) {
	for _, v := range n {
		delete(m, v)
	}
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

func (s Set[T]) Exist(v T) bool {
	_, ok := s[v]
	return ok
}

func (s Set[T]) Remove(v T) {
	delete(s, v)
}

func (s Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s))
	for k := range s {
		slice = append(slice, k)
	}
	return slice
}

func Zip[T comparable, U any](k []T, v []U) map[T]U {
	if len(k) > len(v) {
		m := make(map[T]U, len(k))
		for i := range k {
			m[k[i]] = v[i]
		}
		return m
	} else {
		m := make(map[T]U, len(v))
		for i := range v {
			m[k[i]] = v[i]
		}
		return m
	}
}

func FromSliceFunc[T, U comparable](slice []T, f func(int, T) U) Set[U] {
	s := make(Set[U], len(slice))
	for i, v := range slice {
		s.Add(f(i, v))
	}
	return s
}

func FromSlice[T comparable](slice []T) Set[T] {
	s := make(Set[T], len(slice))
	for _, v := range slice {
		s.Add(v)
	}
	return s
}

func MapMul[T comparable, U Integer](m map[T]U, mul U) {
	for k := range m {
		m[k] *= mul
	}
}

func GetMaxKeyAndValue[T Integer, U any](m map[T]U) (T, U) {
	maxKey := GetMaxKey(m)
	return maxKey, m[maxKey]
}

func GetMaxKey[T Integer, U any](m map[T]U) T {
	var maxKey T
	for k := range m {
		if k > maxKey {
			maxKey = k
		}
	}
	return maxKey
}
