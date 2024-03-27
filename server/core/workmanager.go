package core

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const MaxTaskInChan = 1000

// 定义接口
// Task 定义任务类型
type Task interface {
	Process()
}

// 定义一个回调函数，用于清理资源
type WorkerCleanF func(workerId int64)

// Worker 接口定义工作者行为
type Worker interface {
	Init(id int64, taskChan chan Task, wg *sync.WaitGroup, f WorkerCleanF)
	Start()
	Stop()
}

// //////////////////////////////////////////////////
// 定义最基础的2个类，用作基类
type BaseTask struct {
	Id int64
}

// Run 实现 Task 接口的 Run 方法
func (t *BaseTask) Process() {
	//fmt.Printf("Task %d is running\n", t.Id)
	time.Sleep(time.Second * 1)
	fmt.Printf("Task %d exits\n", t.Id)

}

type BaseWorker struct {
	Id       int64
	waitGrp  *sync.WaitGroup
	taskChan chan Task
	cleanFun WorkerCleanF //相当于析构函数
	quitChan chan struct{}
}

func (w *BaseWorker) Init(id int64, tc chan Task, wg *sync.WaitGroup, f WorkerCleanF) {
	w.Id = id
	w.waitGrp = wg
	w.taskChan = tc
	w.cleanFun = f
	w.quitChan = make(chan struct{}) // 等待退出
	//fmt.Println("init worker ", w.Id)

}
func (w *BaseWorker) Start() {
	w.waitGrp.Add(1)
	defer func() {
		w.waitGrp.Done()
		if w.cleanFun != nil {
			w.cleanFun(w.Id)
		}

	}()

	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()

	for {
		select {
		case task := <-w.taskChan: // 从taskChan接收任务
			// 执行任务处理逻辑
			//fmt.Printf("Worker %d processing task: %#v\n", w.Id, task)
			// ... 这里添加实际的任务处理代码 ...
			task.Process()

			// 重置计时器
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(30 * time.Second)

		case <-timer.C: // 超时处理
			fmt.Printf("Worker %d timed out, exiting...\n", w.Id)
			return

		case <-w.quitChan: // 收到退出信号，结束goroutine
			fmt.Printf("Worker %d received quit signal, exiting...\n", w.Id)
			return
		}
	}

}

func (w *BaseWorker) Stop() {
	close(w.quitChan)
}

func NewBaseWorker() *BaseWorker {
	return &BaseWorker{}
}

// /////////////////////////////////////////////////////
// Manager 定义任务管理器
type Manager[T Task, W Worker] struct {
	workers       map[int64]W
	maxWorkers    int64
	workerCounter int64
	taskChan      chan Task
	lock          sync.Mutex
	wg            sync.WaitGroup

	newWorkerFunc func() W
	exiting       int32
	workerIdSeq   int64 // 可以使用其他算法，这里先用INT64流水号，方便观察，
}

// NewManager 创建新的任务管理器
func NewManager[T Task, W Worker](max int64, newWorkerF func() W) *Manager[T, W] {
	return &Manager[T, W]{
		workers:       make(map[int64]W),
		maxWorkers:    max,
		workerCounter: 0,
		taskChan:      make(chan Task, MaxTaskInChan),
		newWorkerFunc: newWorkerF,
		exiting:       0,
		workerIdSeq:   0,
	}
}

// 清理工作者
func (m *Manager[T, W]) removeWorker(id int64) {
	atomic.AddInt64(&m.workerCounter, -1)
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.workers, id)
}

// AddWorker 添加工作者到管理器
func (m *Manager[T, W]) addWorker(id int64, worker W) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.workers[id] = worker
	atomic.AddInt64(&m.workerCounter, 1) // 计数器原子加1

	go worker.Start()
}

// Wait 等待所有工作者完成任务
func (m *Manager[T, W]) Wait() {
	m.wg.Wait()
}

func (m *Manager[T, W]) StopAll() {
	atomic.StoreInt32(&m.exiting, 1)
	var workers []W

	{
		m.lock.Lock()
		defer m.lock.Unlock()

		workers = make([]W, 0, len(m.workers))
		// 遍历map并将键拷贝到切片中
		for _, value := range m.workers {
			workers = append(workers, value)
		}
	}
	// 上一个加锁需要用代码块包围
	for _, w := range workers {
		w.Stop() // 协程的调用清理函数时候会加锁
	}

}

// 添加任务到队列中
// 工作者少于2个，队列中多于5个任务，则新建立几个
func (m *Manager[T, W]) AddTask(task T) error {

	if atomic.LoadInt32(&m.exiting) == 1 {
		return errors.New("exiting now, can't add task. ")
	}
	// 写入管道
	m.taskChan <- task

	current := atomic.LoadInt64(&m.workerCounter)
	if (current < 2) || (current < m.maxWorkers && len(m.taskChan) > 5) {
		if m.newWorkerFunc == nil {
			return errors.New("create worker function is null")
		}
		worker := m.newWorkerFunc()
		//id := Globals.snow.GenerateID()
		id := atomic.AddInt64(&m.workerIdSeq, 1)
		worker.Init(id, m.taskChan, &m.wg, m.removeWorker)
		m.addWorker(id, worker)
	}
	return nil
}

//////////////////////////////////////////////////////////////////
