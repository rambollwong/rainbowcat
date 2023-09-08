package task

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrRegistered = errors.New("registered task type")
)

var _ Monitor = (*TasksMonitor)(nil)

type Monitor interface {
	Start() error
	Stop() error
	SetDataStore(store DataStore)
	Registered(taskType Type) bool
	RegisterTimerForTasks(triggerTime time.Time, taskType Type, handler Handler) error
	RegisterTickerForTasks(interval time.Duration, taskType Type, handler Handler) error
}

type TimerTask struct {
	tm          *TasksMonitor
	timer       *time.Timer
	taskType    Type
	triggerTime time.Time
	handler     Handler
}

func (t *TimerTask) Run() {
	interval := t.triggerTime.Sub(time.Now())
	if interval < 0 {
		return
	}
	t.timer = time.NewTimer(interval)
	select {
	case <-t.timer.C:
		t.handler(t.tm.dataStore.GetData(t.taskType))
	case <-t.tm.ctx.Done():
		if !t.timer.Stop() {
			<-t.timer.C
		}
		return
	case <-t.tm.exitC:
		if !t.timer.Stop() {
			<-t.timer.C
		}
		return
	}
}

type TickerTask struct {
	tm       *TasksMonitor
	ticker   *time.Ticker
	taskType Type
	interval time.Duration
	handler  Handler
}

func (t *TickerTask) Run() {
	t.ticker = time.NewTicker(t.interval)
	for {
		select {
		case <-t.ticker.C:
			t.handler(t.tm.dataStore.GetData(t.taskType))
		case <-t.tm.ctx.Done():
			t.ticker.Stop()
			return
		case <-t.tm.exitC:
			t.ticker.Stop()
			return
		}
	}
}

type TasksMonitor struct {
	ctx       context.Context
	dataStore DataStore

	mu        sync.RWMutex
	once      sync.Once
	running   bool
	timerMap  map[Type]*TimerTask
	tickerMap map[Type]*TickerTask

	exitC chan struct{}
}

func (t *TasksMonitor) Start() error {
	var err error
	t.once.Do(func() {
		t.exitC = make(chan struct{})

		t.mu.Lock()
		defer t.mu.Unlock()
		for _, task := range t.timerMap {
			go task.Run()
		}
		for _, task := range t.tickerMap {
			go task.Run()
		}
		t.running = true
	})
	return err
}

func (t *TasksMonitor) Stop() error {
	defer func() {
		t.once = sync.Once{}
	}()
	t.mu.Lock()
	defer t.mu.Unlock()
	close(t.exitC)
	t.running = false
	return nil
}

func (t *TasksMonitor) SetDataStore(store DataStore) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.dataStore = store
}

func (t *TasksMonitor) Registered(taskType Type) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.timerMap[taskType]
	if !ok {
		_, ok = t.tickerMap[taskType]
	}
	return ok
}

func (t *TasksMonitor) RegisterTimerForTasks(triggerTime time.Time, taskType Type, handler Handler) error {
	if t.Registered(taskType) {
		return ErrRegistered
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	newTimer := &TimerTask{
		tm:          t,
		taskType:    taskType,
		triggerTime: triggerTime,
		handler:     handler,
	}
	t.timerMap[taskType] = newTimer
	if t.running {
		go newTimer.Run()
	}
	return nil
}

func (t *TasksMonitor) RegisterTickerForTasks(interval time.Duration, taskType Type, handler Handler) error {
	if t.Registered(taskType) {
		return ErrRegistered
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	newTicker := &TickerTask{
		tm:       t,
		taskType: taskType,
		interval: interval,
		handler:  handler,
	}
	t.tickerMap[taskType] = newTicker
	if t.running {
		go newTicker.Run()
	}
	return nil
}
