package main

import (
	"context"
	"demo-task/payload"
	"demo-task/periodic"
	"demo-task/utils"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

func serverErrorHandler(ctx context.Context, task *asynq.Task, err error) {
	// 不能直接判断err是否为nil 因为没有注册任务之前 任务的payload都是空的，所以下面会保错误: unexpected end of JSON input
	if len(task.Payload()) == 0 {
		return
	}
	if err != nil {
		log.Fatalf("任务服务运行异常：%s 错误信息: %v", task.Type(), err)
	}
}

func main() {
	redisClientOpt := asynq.RedisClientOpt{Addr: "localhost:6379"}
	srv := asynq.NewServer(
		redisClientOpt,
		asynq.Config{
			Concurrency:  10,
			ErrorHandler: asynq.ErrorHandlerFunc(serverErrorHandler),
		},
	)

	mux := asynq.NewServeMux()
	// 欢迎邮件任务 具体执行
	mux.HandleFunc(utils.TASK_EMAIL_WELCOME, sendWelcomeEmail)
	// 提醒邮件任务 具体执行
	mux.HandleFunc(utils.TASK_EMAIL_REMINDER, sendReminderEmail)
	// 周期邮件任务 具体执行
	mux.HandleFunc(utils.TASK_EMAIL_PERIODIC, sendPeriodicEmail)

	// 周期任务（动态）
	provider := periodic.FileBasedConfigProvider{Filename: "periodic/periodic_task_config.yaml"}
	manager, err := periodic.GetPeriodicTaskManager(&provider, redisClientOpt)
	if err != nil {
		log.Fatal(err)
	}

	if err = manager.Start(); err != nil {
		log.Fatal(err)
	}

	// 周期任务（静态）
	scheduler := asynq.NewScheduler(redisClientOpt, &asynq.SchedulerOpts{
		EnqueueErrorHandler: utils.HandleEnqueueError,
	})
	payloadSchedule, err := json.Marshal(payload.EmailTaskPayload{Username: "孙行者"})
	if err != nil {
		log.Fatal(err)
	}
	taskSchedule := asynq.NewTask(utils.TASK_EMAIL_WELCOME, payloadSchedule)
	entryID, err := scheduler.Register("@every 30s", taskSchedule)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)
	if err = scheduler.Start(); err != nil {
		log.Fatal(err)
	}

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}

func sendWelcomeEmail(ctx context.Context, t *asynq.Task) error {
	var payload payload.EmailTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	log.Printf(" [*] 任务执行-欢迎邮件发送 %s", payload.Username)
	return nil
}

func sendReminderEmail(ctx context.Context, t *asynq.Task) error {
	var payload payload.EmailTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	log.Printf(" [*] 任务执行-提醒邮件发送 %s", payload.Username)
	return nil
}

func sendPeriodicEmail(ctx context.Context, t *asynq.Task) error {
	var payload payload.EmailTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	log.Printf(" [*] 周期任务执行-周期邮件发送 %s", payload.Username)
	return nil
}
