# 钩子函数设计

本质上来说，是在执行命令前后，额外执行一些逻辑。中间件模式

# 结构设计

```go

// 对外提供的结构体
type hooksMixin struct {
	slice   []Hook
	initial hooks // 初始 这块属于系统内置的，内置的是用来处理 redis 命令的
	current hooks // 当前 外层有套了用户自定义的一些钩子函数
}

// 隐藏在内部的 hooks 结构体，代表了具体的钩子能力
type hooks struct {
	dial       DialHook            // 当创建网络连接时调用的hook
	process    ProcessHook         // 执行命令时调用的hook
	pipeline   ProcessPipelineHook // 执行管道命令时调用的hook
	txPipeline ProcessPipelineHook // 执行事务管道命令时调用的hook
}

```

## 初始化逻辑

```go
// 这个 initHook 函数，是在配置系统实现的 redis 逻辑函数，作为钩子函数的初始值

func (c *Client) init() {
	c.cmdable = c.Process
	c.initHooks(hooks{
		dial:       c.baseClient.dial, // tcp 链接
		process:    c.baseClient.process, // 发送 redis 命令
		pipeline:   c.baseClient.processPipeline, 
		txPipeline: c.baseClient.processTxPipeline,
	})
}

```

## 添加钩子函数