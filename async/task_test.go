package async

import (
	"sync"
	"testing"
)

func TestBatchTask(t *testing.T) {
	cnt := 125923_421
	bs := 7274
	invCnt := 0
	resultWg := sync.WaitGroup{}
	resultWg.Add(cnt)
	bt := NewBatchTask(bs, func(l []int) {
		invCnt += 1
		resultWg.Add(-len(l))
	})
	bt.Watch()

	sendWg := sync.WaitGroup{}
	sendWg.Add(cnt)
	for i := 0; i < cnt; i++ {
		// go func(i int) { // 异步并发send
		bt.Add(i)
		sendWg.Done()
		// }(i)
	}

	sendWg.Wait()
	bt.Stop()

	resultWg.Wait()
	t.Logf("times:%d estimate:%d metric:%v", invCnt, cnt/bs, bt.Metric())
}

func TestClean(t *testing.T) {
	arr := make([]int, 2, 5)
	arr[0] = 1
	arr[1] = 2
	t.Logf("len:%d cap:%d %v", len(arr), cap(arr), arr)
	clear(arr)
	t.Logf("len:%d cap:%d %v", len(arr), cap(arr), arr)
	arr = arr[:0]
	t.Logf("len:%d cap:%d %v", len(arr), cap(arr), arr)
}
