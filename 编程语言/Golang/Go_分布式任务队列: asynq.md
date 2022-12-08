## Asynq简单、可靠、高效的分布式任务队列

> 今天介绍 一下在写项目中使用到的一个简单高效的任务队列库。

Asynq 是一个 Go 库，用于排队任务并与 worker 异步处理它们。它由 Redis 提供支持，旨在实现可扩展且易于上手。

### 1. 概述

Asynq 工作原理的高级概述：

- 客户端将任务放入队列
- 服务器从队列中拉取任务并为每个任务启动一个工作协程
- 任务由多个worker同时处理

任务队列用作跨多台机器分配工作的机制。一个系统可以由多个工作服务器和代理组成，使其具有高可用性和水平扩展特性。

示例用例：

![Example use case](https://user-images.githubusercontent.com/11155743/116358505-656f5f80-a806-11eb-9c16-94e49dab0f99.jpg)

### 2. 特性

Asynq有很多易用的特性，下面简单列举几个：

- 任务调度
- 保证任务至少执行一次
- 失败任务重试
- 自动恢复
- 优先级队列
- 使用唯一选项对任务进行重复数据删除
- 周期性任务
- 可视化的管理界面
- 支持Redis集群和哨兵模式
- 。。。

### 3. 示例

#### 3.1 工具类

定义一个工具类方便使用：`utils.go`

```go
package utils

import (
	"log"

	"github.com/hibiken/asynq"
)
// 任务的task_types，可以认为是任务名称 用来任务执行中的关联
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
```

定义一个任务参数类：`payload.go`

```go
package payload

type EmailTaskPayload struct {
	Username string `yaml:"username"`
}
```

#### 3.2 定义任务

在 `asynq` 中，一个工作单元被封装在一个名为 `Task` 的类型中，它在概念上有两个字段：`Type`和 `Payload`。

```go
// Type 是一个字符串值，表示任务的类型。
func (t *Task) Type() string

// Payload是任务执行所需要的数据。
func (t *Task) Payload() []byte
```

创建任务：

```go
weTask := asynq.NewTask(utils.TASK_EMAIL_WELCOME, payload) // 欢迎邮件发送任务
reTask := asynq.NewTask(utils.TASK_EMAIL_REMINDER, payload) // 提醒邮件发送任务
```

#### 3.3 任务入队

1. 即时排列任务入队：如果没有设置`ProcessIn`或者`ProcessAt`，任务将立即进入待办队列，如果空闲，则会立即得到执行。

```go
// 首先定义一个client
client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})

// 排列入队
// 如果任务入队成功，Enqueue 返回 TaskInfo 和 nil 错误，否则返回非 nil 错误。
info,err := client.Enqueue(weTask)
```

2. 使用 `ProcessIn` 或 `ProcessAt` 选项来安排将来要处理的任务。

```go
// 24 小时后处理任务。
info, err = client.Enqueue(reTask, asynq.ProcessIn(24*time.Hour))

// 现在时间+24小时后 处理任务
client.Enqueue(weTask, asynq.ProcessAt(time.Now().Add(24*time.Hour)))
```

`ProcessIn`接受一个time.Duration类型参数，`ProcessAt`接受一个time.Time类型参数。

#### 3.4 任务处理

**创建任务处理的`asynq.Server`。**

```go
// NewServer 接受一个RedisConnOpt和Config参数。
// Config用于调整服务器的任务处理行为。
srv := asynq.NewServer(
  asynq.RedisClientOpt{Addr: "localhost:6379"},
  asynq.Config{
    Concurrency: 10,  // Concurrency表示最大并发处理任务数。
  },
)
```

- 启动方式1: 开始任务处理并阻塞

```go
// Run 开始任务处理并阻塞，直到接收到退出程序的操作系统信号。 
// 一旦它收到一个信号，它就会优雅地关闭所有活跃的工作人员和其他 goroutines 来处理任务。
if err := srv.Run(handler); err != nil {
  log.Fatal(err)
}
```

- 启动方式2: 开始任务处理并非阻塞

```go
// Start 启动工作服务器。 服务器启动后，它会从队列中取出任务并为每个任务启动一个工作协程，然后调用 Handler 来处理它。
// 任务由workers并发处理，最多达到 Config.Concurrency 中指定的并发数。
if err := srv.Start(handler); err != nil {
  log.Fatal(err)
}
```

> 这种方式，如果你的`main`函数没有自己增加阻塞退出的方法，任务处理会出现问题。因为`main`函数处理完后，程序会立即退出。

- 启动方式3:  ServeMux

```go
// 就像“net/http”包中的 ServeMux 一样，可以通过调用 Handle 或 HandleFunc 来注册处理程序。
// ServeMux 满足 Handler 接口，这样你就可以将它传递给 (*Server).Run或者Start。
mux := asynq.NewServeMux()
// 欢迎邮件任务 具体执行
mux.HandleFunc(utils.TASK_EMAIL_WELCOME, sendWelcomeEmail)
```

**下面介绍 一下上面用到的`handler`。**

要求`handler`必须要具有`ProcessTask`方法。

```go
type Handler interface {
    // 如果任务处理成功，ProcessTask 应该返回 nil。
    // 如果 ProcessTask 返回一个非零错误或panics，任务将在稍后重试。
    ProcessTask(context.Context, *Task) error
}
```

- 示例1:

```go
// 任务处理
func handler(ctx context.Context, t *asynq.Task) error {
    switch t.Type() {
    case TASK_EMAIL_WELCOME :
        var p EmailTaskPayload
        if err := json.Unmarshal(t.Payload(), &p); err != nil {
            return err
        }
    case TASK_EMAIL_REMINDER:
        var p EmailTaskPayload
        if err := json.Unmarshal(t.Payload(), &p); err != nil {
            return err
        }
    default:
        return fmt.Errorf("unexpected task type: %s", t.Type())
    }
    return nil
}

func main() {
    srv := asynq.NewServer(
        asynq.RedisClientOpt{Addr: "localhost:6379"},
        asynq.Config{Concurrency: 10},
    )

    // 需要使用asynq.HandlerFunc适配
    if err := srv.Run(asynq.HandlerFunc(handler)); err != nil {
        log.Fatal(err)
    }
}
```

 我们可以继续向这个处理函数添加 switch case，但在实际的应用程序中，在单独的函数中为每个 case 定义逻辑很方便。

- 示例2:

```go
func main() {
  mux := asynq.NewServeMux()
  // 欢迎邮件任务 具体执行
  mux.HandleFunc(utils.TASK_EMAIL_WELCOME, sendWelcomeEmail)
  // 提醒邮件任务 具体执行
  mux.HandleFunc(utils.TASK_EMAIL_REMINDER, sendReminderEmail)
  // 周期邮件任务 具体执行
  mux.HandleFunc(utils.TASK_EMAIL_PERIODIC, sendPeriodicEmail)
  
  srv := asynq.NewServer(
    asynq.RedisClientOpt{Addr: "localhost:6379"},
    asynq.Config{Concurrency: 10},
  )
  
  // 不需要适配，因为ServeMux实现了Handler接口
  if err := srv.Run(mux); err != nil {
    log.Fatal(err)
  }
}
```

该代码看起来更加有条理。

### 3.5 定时/周期任务

定时任务的创建和普通任务一样，只不过是注册方式不一样。

- 静态

```go
// 创建任务
task := asynq.NewTask(TASK_EMAIL_WELCOME, payload)

// 创建调度器
scheduler := asynq.NewScheduler(redisConnOpt, &asynq.SchedulerOpts{})

// 注册任务
// 可以使用 cron 规范字符串来指定计划。
entryID, err := scheduler.Register("* * * * *", task) // 每分钟执行一次
// 也可以使用"@every <duration>"语法指定间隔
entryID, err = scheduler.Register("@every 30s", task) // 每30s执行一次
// 可以在注册任务的同时，指定配置项
entryID, err = scheduler.Register("@every 24h", task, asynq.Queue("myqueue")) // 每24小时执行一次 队列名字“myqueue”。

// 运行调度器 并阻塞
if err := scheduler.Run(); err != nil {
    log.Fatal(err)
}
```

- 动态

如果要动态添加和删除周期性任务（即不重新启动Scheduler进程），可以使用`PeriodicTaskManager`。`PeriodTaskManager`使用`PeriodTaskConfigProvider`定期获取当前定期任务配置，并将调度器的条目与当前配置同步。

可以将周期性任务配置存储在数据库或本地文件中，并更新此配置源以动态添加和删除周期性任务。也可以轻松修改示例以使用数据库或其他配置源。

使用YAML文件进行说明。

1. 定义任务信息

```yaml
configs:
  - cronspec: "* * * * *"
    task_type: "task:email:periodic"
    payload:
    	username: "孙悟空"
```

2. 读取文件内容

```go
type FileBasedConfigProvider struct {
     filename string
}

// 必须要实现的方法(接口：PeriodicTaskConfigProvider)，读取所有配置项
func (p *FileBasedConfigProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
    data, err := os.ReadFile(p.filename)
    if err != nil {
        return nil, err
    }
    var c PeriodicTaskConfigContainer
    if err := yaml.Unmarshal(data, &c); err != nil {
        return nil, err
    }
    var configs []*asynq.PeriodicTaskConfig
    for _, cfg := range c.Configs {
         configs = append(configs, &asynq.PeriodicTaskConfig{Cronspec: cfg.Cronspec, Task: asynq.NewTask(cfg.TaskType, nil)})
    }
    return configs, nil
}

type PeriodicTaskConfigContainer struct {
    Configs []*Config `yaml:"configs"`
}

type Config struct {
    Cronspec string `yaml:"cronspec"`
    TaskType string `yaml:"task_type"`
}
```

3. 创建任务管理器

```go
func main() {
	  // 指定文件名称
    provider := &FileBasedConfigProvider{filename: "./periodic_task_config.yml"}

    mgr, err := asynq.NewPeriodicTaskManager(
        asynq.PeriodicTaskManagerOpts{
            RedisConnOpt:               asynq.RedisClientOpt{Addr: "localhost:6379"},
            PeriodicTaskConfigProvider: provider,         // 配置源的接口
            SyncInterval:               10 * time.Second, // 指定同步发生的频率（同步配置源）
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

4. 运行任务管理器

```go
// 运行 并阻塞
if err := mgr.Run(); err != nil {
  log.Fatal(err)
}

```

当它运行时，尝试更改配置文件。应该会看到表明添加或删除了新配置的日志消息。

### 4. 错误处理

#### 4.1 定期任务

```go
func handleEnqueueError(task *asynq.Task, opts []asynq.Option, err error) {
    // 错误处理逻辑
}

scheduler := asynq.NewScheduler(
    redisConnOpt, 
    &asynq.SchedulerOpts{
        EnqueueErrorHandler: handleEnqueueError,
    },
)
```

#### 4.2 定期任务（动态）

```go
func handleEnqueueError(task *asynq.Task, opts []asynq.Option, err error) {
  // 错误处理逻辑
}
manager := asynq.NewPeriodicTaskManager(
  asynq.PeriodicTaskManagerOpts{
    RedisConnOpt:               redisOpt,
    PeriodicTaskConfigProvider: provider,
    SyncInterval:               10 * time.Second, // 和配置文件同步的时间，配置文件如果修改，可以及时运用到程序
    SchedulerOpts: &asynq.SchedulerOpts{
      EnqueueErrorHandler: handleEnqueueError,
    },
  },
)
```

#### 4.3 任务处理服务

```go
func serverErrorHandler(ctx context.Context, task *asynq.Task, err error) {
  // 错误处理
}
// ErrorHandlerFunc 类型是一个适配器，允许将普通函数用作 ErrorHandler。
// 如果 f 是具有适当签名的函数，则 ErrorHandlerFunc(f) 是调用 f 的 ErrorHandler。
srv := asynq.NewServer(
  redisClientOpt,
  asynq.Config{
    Concurrency:  10,
    ErrorHandler: asynq.ErrorHandlerFunc(serverErrorHandler),
  },
)
```

### 5. 注意事项

- 该库目前正在进行大量开发，API 经常发生重大变化。

  > 官方说明：当前主要版本为零 (v0.x.x) 以适应快速开发和快速迭代，同时获得用户的早期反馈（感谢对 API 的反馈！）。在 v1.0.0 发布之前，公共 API 可以在没有主要版本更新的情况下更改。

- `(*Server).Run`和`(*Server).Start`一定要分清。一个是阻塞的，一个是非阻塞的。
- **定时任务使用的是[cron](https://github.com/robfig/cron/)库，cron 规范字符串遵循该库的规范。只支持五位的cron，从分钟开始。**
- 错误处理也要注意，如果没有任务放入队列的话，直接判断是否有错误，则会出现：`unexpected end of JSON input`异常。

## 结语

总而言之，Asynq是一个非常易于使用库。其他使用信息可以在[GitHub](https://github.com/hibiken/asynq)上查看。
