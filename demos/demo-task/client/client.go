package main

import (
	"demo-task/payload"
	"demo-task/utils"
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})

	payload, err := json.Marshal(payload.EmailTaskPayload{Username: "孙悟空"})
	if err != nil {
		log.Fatal(err)
	}

	weTask := asynq.NewTask(utils.TASK_EMAIL_WELCOME, payload)
	reTask := asynq.NewTask(utils.TASK_EMAIL_REMINDER, payload)

	// 加入队列，并立即执行任务
	info, err := client.Enqueue(weTask)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [*] Successfully enqueued task: %+v", info)

	// 加入队列，60s以后执行
	info, err = client.Enqueue(reTask, asynq.ProcessIn(60*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [*] Successfully enqueued task: %+v", info)

}
