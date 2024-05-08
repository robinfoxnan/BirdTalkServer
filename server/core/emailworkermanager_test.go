package core

import (
	"birdtalk/server/email"
	"strconv"
	"testing"
	"time"
)

func TestEmailWorkers(t *testing.T) {
	manager := NewEmailWorkerManager(20)

	// 添加示例工作者到管理器

	// 添加示例任务到管理器
	go func() {
		for i := 0; i < 10; i++ {
			var t = &EmailTask{
				BaseTask: BaseTask{Id: int64(i)},
				data: &email.EmailData{
					HostUrl: "http://birdtalk.com",
					Code:    "12345 " + strconv.Itoa(i),
					Session: "123333333333",
					Server:  "1",
					Email:   "robin-fox@sohu.com",
				},
			}

			manager.AddTask(t)
		}
	}()

	//time.Sleep(time.Second * 50)

	time.Sleep(time.Second * 25)
	manager.StopAll()
	// 等待所有工作者完成任务
	manager.Wait()
}
