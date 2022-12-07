## asynq入门示例

示例介绍：

1. client加入任务到队列，可以指定任务执行的类型，是立即执行或者未来某个时间执行。
2. worker是任务的执行者，根据加入的任务类型进行执行。
3. 周期任务只会执行一次，暂时没有找到原因。。。。。

更详细的资料请查看[官方介绍](https://github.com/hibiken/asynq/wiki/Getting-Started)

> 周期任务的执行依赖：[cron](https://github.com/robfig/cron)库，所以cron表达式要根据此库的规范写，只能有五项，从分钟开始，没有秒。要使用秒可以用`@every 10s`语法。