package utils

import (
	"log"

	"github.com/hibiken/asynq"
)

const (
	TASK_EMAIL_WELCOME  = "task:email:welcome"
	TASK_EMAIL_REMINDER = "task:email:reminder"
	TASK_EMAIL_PERIODIC = "task:email:periodic"
)

// 错误处理
func HandleEnqueueError(task *asynq.Task, opts []asynq.Option, err error) {
	if err != nil {
		log.Fatal(task.Type(), err)
	}
}
