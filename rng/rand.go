package rng

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"math/rand"
)

// l < h
func Range(l, h int) int {
	return rand.Intn(h-l) + l
}

// [min,max]
func BetweenClose(i interface {
	GetMin() uint32
	GetMax() uint32
}) uint32 {
	l, r := i.GetMin(), i.GetMax()
	if l == r {
		return l
	}
	return rand.Uint32()%(r-l+1) + l
}

// [min,max)
func BetweenOpen(i interface {
	GetMin() uint32
	GetMax() uint32
}) uint32 {
	if i.GetMin() == i.GetMax() {
		return i.GetMin()
	}
	return rand.Uint32()%(i.GetMax()-i.GetMin()) + i.GetMin()
}

// [0,n)
func Int(n int) int {
	return rand.Intn(n)
}

func In[T any](list []T) T {
	return list[Int(len(list))]
}

// f in [0.0,1.0)
func Exceed(f float64) bool {
	return f > rand.Float64()
}

func WeightRandF[T any](list []T, p func(T) uint64) T {
	if len(list) == 1 {
		return list[0]
	}
	var sum uint64
	for i := range list {
		sum += p(list[i])
	}
	rs := uint64(rand.Int63n(int64(sum)))
	for i := range list {
		ii := p(list[i])
		if rs < ii {
			return list[i]
		}
		rs -= ii
	}
	panic(fmt.Sprintf("list: %v rand: %d", list, rs))
}

func WeightRandI[T constraints.Integer](list []T) int {
	var sum uint64
	for _, v := range list {
		sum += uint64(v)
	}
	rs := uint64(rand.Int63n(int64(sum)))
	for i, v := range list {
		if rs < uint64(v) {
			return i
		}
		rs -= uint64(v)
	}
	panic(fmt.Sprintf("list: %v rand: %d", list, rs))
}
