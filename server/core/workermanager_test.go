package core

import (
	"fmt"
	"testing"
	"time"
)

type CustomTask struct {
	BaseTask
	AdditionalInfo string
}

// 实现 Task 接口的 Run 方法，调用父类 BaseTask 的 Process 方法
func (t *CustomTask) Process(w Worker) {
	//fmt.Printf("CustomTask with additional info '%s' is running\n", t.AdditionalInfo)
	// 调用父类的 Process 方法
	//t.BaseTask.Process()
	fmt.Println("Custom task exit", t.BaseTask.Id)
}

func TestWorkers(t *testing.T) {
	manager := NewManager[Task, *BaseWorker](20, NewBaseWorker)

	// 添加示例工作者到管理器

	// 添加示例任务到管理器
	go func() {
		for i := 0; i < 10; i++ {
			var t = &BaseTask{Id: int64(i)}
			manager.AddTask(t)
		}
		for i := 10; i < 16; i++ {
			var t = &CustomTask{BaseTask: BaseTask{Id: int64(i)}, AdditionalInfo: "Custom Info"}
			manager.AddTask(t)
		}
	}()

	time.Sleep(time.Second * 50)
	go func() {
		for i := 16; i < 26; i++ {
			var t = &BaseTask{Id: int64(i)}
			manager.AddTask(t)
		}
	}()

	time.Sleep(time.Second * 15)
	manager.StopAll()
	// 等待所有工作者完成任务
	manager.Wait()
}
