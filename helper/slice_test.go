package helper

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func TestBSearch(t *testing.T) {
	N := 1000
	list := make([]int, 0, N)
	for i := 0; i < N; i++ {
		list = append(list, i)
	}
	for i := 0; i < N; i++ {
		v := rand.Int() % (N * 2)
		x := BSearch(list[:i], func(ele int) int {
			return v - ele
		})
		if x != -1 {
			assert.Equal(t, v, list[x], "rand:%d index:%d", v, x)
		} else {
			assert.NotContains(t, list[:i], v, "rand:%d index:%d", v, x)
		}
	}

	list = []int{1, 1, 1, 1, 1, 1}
	assert.Equal(t, 3, BSearch(list, func(ele int) int {
		return 1 - ele
	}))
	list = []int{1, 1, 1, 1, 1}
	assert.Equal(t, 2, BSearch(list, func(ele int) int {
		return 1 - ele
	}))
	list = []int{1, 1, 2, 3, 4, 5}
	assert.Equal(t, 1, BSearch(list, func(ele int) int {
		return 1 - ele
	}))
}

func BenchmarkArr(b *testing.B) {
	ss := b2s3(s2b3(strings.Repeat("A", 1024*1024*1024)))
	_ = ss
	b.ReportAllocs()
}

func BenchmarkCopy(b *testing.B) {
	ss := b2s4(s2b4(strings.Repeat("A", 1024*1024*1024)))
	_ = ss
	b.ReportAllocs()
}

func BenchmarkRef(b *testing.B) {
	ss := b2s2(s2b2(strings.Repeat("A", 1024*1024*1024)))
	_ = ss
	b.ReportAllocs()
}

func BenchmarkUnsafe(b *testing.B) {
	ss := b2s1(s2b1(strings.Repeat("A", 1024*1024*1024)))
	_ = ss
	b.ReportAllocs()
}

func s2b1(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func b2s1(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func b2s2(b []byte) string {
	ss := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	yy := (*reflect.StringHeader)(unsafe.Pointer(ss))
	return *(*string)(unsafe.Pointer(yy))
}

func b2s3(b []byte) string {
	y := (*[2]uintptr)(unsafe.Pointer(&b))
	h := [3]uintptr{y[0], y[1], y[1]}
	return *(*string)(unsafe.Pointer(&h))
}

func b2s4(b []byte) string {
	return string(b)
}

func s2b2(s string) []byte {
	ss := (*reflect.StringHeader)(unsafe.Pointer(&s))
	yy := (*reflect.SliceHeader)(unsafe.Pointer(ss))
	return *(*[]byte)(unsafe.Pointer(yy))
}

func s2b3(s string) []byte {
	y := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{y[0], y[1], y[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func s2b4(s string) []byte {
	return []byte(s)
}

func TestLooseEqual(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		assert.True(t, LooseEqual(10, 11, 2))
		assert.True(t, LooseEqual(0, 0, 1))
		assert.False(t, LooseEqual(5, 11, 1))
	})
	t.Run("overflow", func(t *testing.T) {
		assert.False(t, LooseEqual(uint8(254), 255, 3))
		assert.False(t, LooseEqual(uint8(2), 1, 3))

		assert.False(t, LooseEqual(int8(126), 127, 3))
		assert.False(t, LooseEqual(int8(-126), -127, 3))
	})
}
