package async

import (
	"fmt"
	"sync"
	"time"
)

type BatchTask[T any] struct {
	buf  []T
	bs   int
	ch   chan []T
	do   func([]T)
	stop chan struct{}

	p         sync.Pool
	mu        sync.Mutex
	startOnce sync.Once

	byChan      uint64
	chanSend    uint64
	byChanTime  time.Duration
	byTimer     uint64
	timerSend   uint64
	byTimerTime time.Duration
}

func NewBatchTask[T any](bufferLen int, f func([]T)) *BatchTask[T] {
	return &BatchTask[T]{
		buf:  make([]T, 0, bufferLen),
		bs:   bufferLen,
		ch:   make(chan []T, 7274),
		do:   f,
		stop: make(chan struct{}),
		p: sync.Pool{
			New: func() any {
				return make([]T, 0, bufferLen)
			},
		},
	}
}

func (b *BatchTask[T]) Metric() string {
	return fmt.Sprintf("chan[times: %d total: %d use: %s] timer[times: %d total: %d use: %s]",
		b.byChan, b.chanSend, b.byChanTime, b.byTimer, b.timerSend, b.byTimerTime)
}

func (b *BatchTask[T]) Stop() {
	b.stop <- struct{}{}
	<-b.stop
	close(b.stop)
}

func (b *BatchTask[T]) Watch() {
	b.startOnce.Do(func() {
		ticker := time.NewTicker(time.Second)
		go func() {
			defer ticker.Stop()
			for {
				select {
				case <-b.stop:
					close(b.ch)
					t0 := time.Now()
					for v := range b.ch { // clean
						b.chanSend += uint64(len(v))
						b.byChan++
						b.do(v)
					}
					if len(b.buf) != 0 {
						b.chanSend += uint64(len(b.buf))
						b.byChan++
						b.do(b.buf)
					}
					b.byChanTime += time.Since(t0)
					b.stop <- struct{}{}
					return
				case <-ticker.C:
					data := b.swap()
					if len(data) != 0 {
						t0 := time.Now()
						b.timerSend += uint64(len(data))
						b.do(data)
						b.byTimerTime += time.Since(t0)
						b.byTimer++
						b.put(data)
					}
				case data := <-b.ch:
					t0 := time.Now()
					b.chanSend += uint64(len(data))
					b.do(data)
					b.byChanTime += time.Since(t0)
					b.byChan++
					b.put(data)
				}
			}
		}()
	})
}

func (b *BatchTask[T]) swap() []T {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.buf) == 0 {
		return nil
	}
	tmp := b.buf
	b.buf = b.get()
	return tmp
}

func (b *BatchTask[T]) Add(v T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buf = append(b.buf, v)
	if len(b.buf) >= b.bs {
		b.ch <- b.buf
		b.buf = b.get()
	}
}

func (b *BatchTask[T]) get() []T {
	return b.p.Get().([]T)
}

func (b *BatchTask[T]) put(v []T) {
	clear(v)
	b.p.Put(v[:0])
}
