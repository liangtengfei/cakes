package periodic

import (
	"demo-task/payload"
	"demo-task/utils"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"gopkg.in/yaml.v2"
)

// func main() {
// 	// 周期任务管理
// 	provider := FileBasedConfigProvider{Filename: "periodic/periodic_task_config.yaml"}

// 	redisClientOpt := asynq.RedisClientOpt{Addr: "localhost:6379"}
// 	manager, err := GetPeriodicTaskManager(&provider, redisClientOpt)
// 	if err != nil {
// 		log.Fatal("周期任务管理 获取失败：", err)
// 	}
// 	if err := manager.Run(); err != nil {
// 		log.Fatal("周期任务管理 运行", err)
// 	}
// }

func GetPeriodicTaskManager(provider asynq.PeriodicTaskConfigProvider, redisOpt asynq.RedisClientOpt) (*asynq.PeriodicTaskManager, error) {
	// provider := &FileBasedConfigProvider{filename: "./periodic_task_config.yaml"}
	return asynq.NewPeriodicTaskManager(
		asynq.PeriodicTaskManagerOpts{
			RedisConnOpt:               redisOpt,
			PeriodicTaskConfigProvider: provider,
			SyncInterval:               10 * time.Second, // 和配置文件同步的时间，配置文件如果修改，可以及时运用到程序
			SchedulerOpts: &asynq.SchedulerOpts{
				EnqueueErrorHandler: utils.HandleEnqueueError,
			},
		},
	)
}

// 本例使用了yaml文件作为配置源，但是，你也可以使用数据库或其他配置源。
type FileBasedConfigProvider struct {
	Filename string
}

// 解析 yaml 文件并返回 PeriodicTaskConfigs 列表
func (p *FileBasedConfigProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
	data, err := os.ReadFile(p.Filename)
	if err != nil {
		return nil, err
	}
	var c PeriodicTaskConfigContainer
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	var configs []*asynq.PeriodicTaskConfig
	for _, cfg := range c.Configs {
		p, err := json.Marshal(cfg.Payload)
		if err != nil {
			p = []byte("{\"username\":\"无用户名\"}")
			log.Println(cfg.TaskType, "未配置用户名")
		}

		configs = append(configs, &asynq.PeriodicTaskConfig{Cronspec: cfg.Cronspec, Task: asynq.NewTask(cfg.TaskType, p)})
	}
	return configs, nil
}

type PeriodicTaskConfigContainer struct {
	Configs []*Config `yaml:"configs"`
}

type Config struct {
	Cronspec string                   `yaml:"cronsepc"`
	TaskType string                   `yaml:"task_type"`
	Payload  payload.EmailTaskPayload `yaml:"payload"`
}
