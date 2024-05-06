package taskutil

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// Job 具体执行函数
type Job func(key string)

// Task 任务结构体
type Task struct {
	// 任务key（任务唯一标识）
	key string
	// 任务执行函数
	job Job
	// 任务执行时间间隔
	executeAt time.Duration
	// 任务执行次数，-1:无限次
	times int
	// 任务所在槽
	slot int
	// 任务所在环
	circle int
}

// TimeWheel 时间轮
type TimeWheel struct {
	// 时间间隔
	interval time.Duration
	// 定时器
	ticker *time.Ticker
	// 当前槽
	currentSlot int
	// 时间轮总槽数
	slotNum int
	// 时间轮槽双向链表数组
	slots []*list.List
	// 时间轮停止信号
	stopCh chan struct{}
	// 任务移除信号
	removeTaskCh chan string
	// 任务添加信号
	addTaskCh chan *Task
	// 任务记录map，将任务key和slots链表数组的element进行映射和管理
	taskRecords sync.Map
	// 同步锁
	mux sync.Mutex
	// 判断时间轮是否运行
	isRun bool
}

// DefaultTimingWheel 默认时间轮
// 返回值：时间轮对象，错误信息
func DefaultTimingWheel() (*TimeWheel, error) {
	return NewTimingWheel(time.Second, 12)
}

// NewTimingWheel 初始化时间轮
// interval:时间间隔
// slotNum:时间轮总槽数
// 返回值：时间轮对象，错误信息
func NewTimingWheel(interval time.Duration, slotNum int) (*TimeWheel, error) {
	// 如果时间间隔小于秒，则创建错误
	if interval < time.Second {
		return nil, errors.New("minimum interval need greater than or equal second")
	}
	// 如果时间轮槽数小于1，则创建错误
	if slotNum <= 0 {
		return nil, errors.New("minimum slotNum need greater than zero")
	}
	// 初始化时间轮
	t := &TimeWheel{
		interval:     interval,
		currentSlot:  0,
		slotNum:      slotNum,
		slots:        make([]*list.List, slotNum),
		stopCh:       make(chan struct{}),
		removeTaskCh: make(chan string),
		addTaskCh:    make(chan *Task),
		isRun:        false,
	}
	// 启动时间轮
	t.start()
	// 返回时间轮对象
	return t, nil
}

// Stop 停止时间轮
func (t *TimeWheel) Stop() {
	// 如果时间轮正在运行
	if t.isRun {
		// 获取同步锁
		t.mux.Lock()
		// 设置时间轮停止
		t.isRun = false
		// 停止定时器
		t.ticker.Stop()
		// 释放同步锁
		t.mux.Unlock()
		// 关闭stopCh
		close(t.stopCh)
	}
}

// AddTask 添加任务
// key:任务唯一标识
// job:任务执行函数
// executeAt:任务执行时间
// times:任务执行次数，-1:无限次
// 返回值：错误信息
func (t *TimeWheel) AddTask(key string, job Job, executeAt time.Duration, times int) error {
	// key不能为空
	if key == "" {
		return errors.New("key is empty")
	}
	// 任务执行时间不能小于时间轮间隔
	if executeAt < t.interval {
		return errors.New("executeAt should be greater than or equal interval")
	}
	// 判断任务执行次数是否合法
	if times < -1 || times == 0 {
		return errors.New("times should be a positive integer, or -1 to indicate unlimited times")
	}
	// 判断是否已添加过任务
	_, ok := t.taskRecords.Load(key)
	// 如果已添加过任务
	if ok {
		return errors.New("key of job already exists")
	}
	// 初始化任务
	task := &Task{
		key:       key,
		job:       job,
		times:     times,
		executeAt: executeAt,
	}
	// 添加任务
	t.addTaskCh <- task
	return nil
}

// RemoveTask 删除任务
// key:任务唯一标识
// 返回值：错误信息
func (t *TimeWheel) RemoveTask(key string) error {
	// key不能为空
	if key == "" {
		return errors.New("key is empty")
	}
	// 删除任务
	t.removeTaskCh <- key
	return nil
}

// start 启动时间轮
func (t *TimeWheel) start() {
	// 判断时间轮是否在运行
	if !t.isRun {
		// 初始化slots数组
		// 根据时间轮总槽数初始化数组结构双向链表list
		for i := 0; i < t.slotNum; i++ {
			t.slots[i] = list.New()
		}
		// 设置定时器时间间隔
		t.ticker = time.NewTicker(t.interval)
		t.mux.Lock()
		t.isRun = true
		// 开启协程执行
		go t.run()
		t.mux.Unlock()
	}
}

// run 时间轮运行函数
func (t *TimeWheel) run() {
	// 循环监听
	for {
		// 监听信号
		select {
		// 监听停止时间轮信号
		case <-t.stopCh:
			return
		// 监听添加任务信号
		case task := <-t.addTaskCh:
			t.addTask(task)
		// 监听删除任务信号
		case key := <-t.removeTaskCh:
			t.removeTask(key)
		// 监听定时器信号
		case <-t.ticker.C:
			t.execute()
		}
	}
}

// addTask 添加任务
// task:任务对象
func (t *TimeWheel) addTask(task *Task) {
	// 计算任务所在槽，任务所在环
	slot, circle := t.calSlotAndCircle(task.executeAt)
	// 任务所在槽
	task.slot = slot
	// 任务所在环
	task.circle = circle
	// 获取指定槽的链表，并将任务添加到链表中
	ele := t.slots[slot].PushBack(task)
	// 将任务key和slots链表数组的element记录到map中
	t.taskRecords.Store(task.key, ele)
}

// calSlotAndCircle 计算任务所在槽和任务所在环
// executeAt:任务执行时间间隔
// 返回值：任务所在槽，任务所在环
func (t *TimeWheel) calSlotAndCircle(executeAt time.Duration) (slot, circle int) {
	// 任务时间间隔 秒
	delay := int(executeAt.Seconds())
	// 当前轮盘表示的时间 秒
	circleTime := len(t.slots) * int(t.interval.Seconds())
	// 计算所在环数
	circle = delay / circleTime
	// 计算时间间隔对应的slot步长
	steps := delay / int(t.interval.Seconds())
	// 计算所在槽数
	slot = (t.currentSlot + steps) % len(t.slots)
	return
}

// removeTask 删除任务
// key:任务唯一标识
func (t *TimeWheel) removeTask(key string) {
	// 从map中获取任务记录
	taskRec, ok := t.taskRecords.Load(key)
	// 如果任务不存在
	if !ok {
		return
	}
	// 获取任务记录对应的链表元素
	ele := taskRec.(*list.Element)
	// 获取任务对象
	task, _ := ele.Value.(*Task)
	// 从链表中删除任务
	t.slots[task.slot].Remove(ele)
	// 从map中删除任务记录
	t.taskRecords.Delete(key)
}

// execute 执行任务
func (t *TimeWheel) execute() {
	// 获取当前槽的链表
	taskList := t.slots[t.currentSlot]
	// 如果链表不为空
	if taskList != nil {
		// 遍历链表，获取链表首元素
		for ele := taskList.Front(); ele != nil; {
			// 获取任务对象
			taskEle, _ := ele.Value.(*Task)
			// 判断任务是否执行 (circle == 0才执行)
			if taskEle.circle > 0 {
				taskEle.circle--
				ele = ele.Next()
				continue
			}
			// 执行任务
			go taskEle.job(taskEle.key)
			// 从map中删除任务记录
			t.taskRecords.Delete(taskEle.key)
			// 从链表中删除任务
			taskList.Remove(ele)
			// 固定次数任务
			// 如果任务执行次数大于1，则重新添加任务
			if taskEle.times > 1 {
				taskEle.times--
				t.addTask(taskEle)
			}
			// 重复任务
			// 如果任务执行次数为-1，则重新添加任务
			if taskEle.times == -1 {
				t.addTask(taskEle)
			}
			ele = ele.Next()
		}
	}
	// 移动到下一个槽
	t.incrCurrentSlot()
}

// incrCurrentSlot 移动到下一个槽
func (t *TimeWheel) incrCurrentSlot() {
	// 移动到下一个槽
	t.currentSlot = (t.currentSlot + 1) % len(t.slots)
}
