package pool

import (
	"sync/atomic"
	"testing"
	"time"
)

type testStruct struct {
	a int
	b []string
	c uint64
	d uint64
	e []*testStruct
	f time.Time
}

func BenchmarkPoolAlloc(b *testing.B) {
	b.ReportAllocs()
	create := atomic.Int32{}
	p := New(func() *testStruct {
		create.Add(1)
		return &testStruct{}
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			obj := p.Get()
			for i := 0; i < 1000; i++ {
				obj.c = obj.d
				obj.c = obj.d*obj.c + obj.c + 999
			}
			p.Put(obj)
		}
	})
	b.ReportMetric(float64(create.Load()), "object(s)")
}
