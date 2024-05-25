package timermanager

import (
	"container/heap"
	"sync"
	"time"
)

type TimerQueue []*Timer

func (pq TimerQueue) Len() int { return len(pq) }

func (pq TimerQueue) Less(i, j int) bool {
	return pq[i].endTime < pq[j].endTime
}

func (pq TimerQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *TimerQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Timer)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *TimerQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type Timer struct {
	TimerOuter       // 执行方法
	id         int64 // 定时器id
	endTime    int64 // 执行时间
	interval   int64 // 间隔时间
	index      int   // 在堆数组中的索引
}

type TimerOuter interface {
	TimeOut()
}

type Manager struct {
	id   int64      // 自增的timerId
	tq   TimerQueue // 定时器
	lock sync.Mutex
}

// 启动调度器
func (m *Manager) scheduler() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			m.RunTimer()
		}
	}
}

// 添加定时器
func (m *Manager) AddTimer(timerOuter TimerOuter, endTime, interval int64) (int64, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.id++
	timer := &Timer{
		id:         m.id,
		TimerOuter: timerOuter,
		endTime:    endTime,
		interval:   interval,
	}

	if m.tq.Len() >= cap(m.tq) {
		l := m.tq.Len() * 2
		queue := make(TimerQueue, 0, l)
		for i := 0; i < m.tq.Len(); i++ {
			queue.Push(m.tq[i])
		}
		m.tq = queue
	}
	// 添加到堆中
	heap.Push(&m.tq, timer)
	return m.id, nil
}

// 删除定时器
func (m *Manager) RemoveTimer(timerId int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, timer := range m.tq {
		if timer.id == timerId {
			heap.Remove(&m.tq, timer.index)
			return
		}
	}
}

// 清除定时器
func (m *Manager) ClearTimer() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.tq = m.tq[:0]
	return
}

// 执行定时器
func (m *Manager) RunTimer() {
	m.lock.Lock()
	defer m.lock.Unlock()
	// 没有需要执行的任务
	if m.tq.Len() <= 0 {
		return
	}

	var execTimer []interface{} // 将要执行的timer
	for m.tq.Len() > 0 {
		// 从堆顶取一个
		tmp := m.tq[0]
		// 时间未到
		if time.Now().Unix() < tmp.endTime {
			break
		}

		timer := heap.Pop(&m.tq).(*Timer)
		execTimer = append(execTimer, timer)

		// 可重复执行
		if timer.interval > 0 {
			timer.endTime += timer.interval
			heap.Push(&m.tq, timer)
		}
	}

	// 执行定时器
	if len(execTimer) > 0 {
		for _, timer := range execTimer {
			timer.(TimerOuter).TimeOut()
		}
	}
}

func NewTimerManager() *Manager {
	m := &Manager{
		tq: make(TimerQueue, 0, 1024),
	}
	go m.scheduler()
	return m
}
