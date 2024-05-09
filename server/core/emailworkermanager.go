package core

import (
	"birdtalk/server/email"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var emailTaskId int64 = 0

type EmailTask struct {
	BaseTask
	data    *email.EmailData
	session *Session
}

func NewEmailTask(sess *Session, emailAddr string, code string) *EmailTask {
	n := atomic.AddInt64(&emailTaskId, 1)

	data := email.EmailData{
		HostUrl: "https://birdtalk.cc",
		Code:    code,
		Session: strconv.FormatInt(sess.Sid, 10),
		Server:  strconv.Itoa(Globals.Config.Server.HostIndex),
		Email:   emailAddr,
	}

	task := EmailTask{
		BaseTask: BaseTask{Id: n},
		session:  sess,
		data:     &data,
	}

	return &task
}

// 发送数据
func SendEmailCode(sess *Session, emailAddr string, code string) {
	t := NewEmailTask(sess, emailAddr, code)
	Globals.emailWorkerManager.AddTask(t)
}

func (t *EmailTask) Process(w Worker) {

	generator, _ := email.NewEmailGenerator(`D:\GBuild\BirdTalkServer\server\emailtemp\email-validation-zh.templ`)

	worker := w.(*EmailWorker)

	subject, txt, _ := generator.GeneratePlainEmail(t.data)
	//fmt.Println(subject, txt)

	client := worker.GetSmtpClient()
	err := client.SendMail([]string{t.data.Email}, subject, txt) // +t.data.Session
	if err != nil {
		//fmt.Println(err, t.Id)
		Globals.Logger.Info("Email sending meet error ", zap.Int64("task id", t.BaseTask.Id), zap.Error(err))
		t.session.NotifyEmailErr()
	}
	//fmt.Println("Email task is finished ", t.BaseTask.Id)
	Globals.Logger.Info("Email task is finished ", zap.Int64("task id", t.BaseTask.Id))
	client.Close()
	time.Sleep(time.Millisecond * 15)
}

//////////////////////////////////////////////////////////////////////////

type EmailWorker struct {
	BaseWorker
	SmtpClient *email.MailValidator
}

// /////////////////////////////////////////
func NewEmailWorker() *EmailWorker {
	return &EmailWorker{
		BaseWorker: BaseWorker{
			Id:       0,
			waitGrp:  nil,
			taskChan: nil,
			cleanFun: nil,
			quitChan: nil,
		},
		SmtpClient: email.NewMailValidator(Globals.Config.Email.SMTPAddr,
			Globals.Config.Email.SMTPPort,
			Globals.Config.Email.SMTPHeloHost,
			Globals.Config.Email.TLSInsecureSkipVerify,
			Globals.Config.Email.UserName,
			Globals.Config.Email.UserPwd),
	}
}

func (w *EmailWorker) GetSmtpClient() *email.MailValidator {
	return w.SmtpClient
}

func (w *EmailWorker) Init(id int64, tc chan Task, wg *sync.WaitGroup, f WorkerCleanF) {
	(&w.BaseWorker).Init(id, tc, wg, f)
	//fmt.Println("init worker ", w.Id)

}
func (w *EmailWorker) Start() {
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
			task.Process(w)

			// 重置计时器
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(30 * time.Second)

		case <-timer.C: // 超时处理
			//fmt.Printf("EmailWorker %d timed out, exiting...\n", w.Id)
			Globals.Logger.Info("Email worker exit because of timeout", zap.Int64("worker id", w.Id))
			return

		case <-w.quitChan: // 收到退出信号，结束goroutine
			//fmt.Printf("EmailWorker %d received quit signal, exiting...\n", w.Id)
			Globals.Logger.Info("Email worker exit because of signal", zap.Int64("worker id", w.Id))
			return
		}
	}
}

func (w *EmailWorker) Stop() {
	w.BaseWorker.Stop()
	w.SmtpClient.Close()
}

// ////////////////////////////////////////////////////////////
func NewEmailWorkerManager(nWorkers int64) *Manager[Task, *EmailWorker] {
	manager := NewManager[Task, *EmailWorker](nWorkers, NewEmailWorker)
	return manager
}

func TestEmailWorkers1() {
	manager := NewEmailWorkerManager(3)

	// 添加示例工作者到管理器

	// 添加示例任务到管理器
	go func() {
		for i := 0; i < 6; i++ {
			var t = &EmailTask{
				BaseTask: BaseTask{Id: int64(i)},
				data: &email.EmailData{
					HostUrl: "http://birdtalk.com",
					Code:    "12345 " + strconv.Itoa(i),
					Session: strconv.Itoa(i),
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
