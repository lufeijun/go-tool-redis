package redis

import (
	"context"

	"github.com/lufeijun/go-tool-redis/redis/command"
)

// 可执行命令函数
type cmdable func(ctx context.Context, cmd command.Cmder) error
type statefulCmdable func(ctx context.Context, cmd command.Cmder) error

type Client struct {
	*baseClient // 基础客户端字段
	cmdable     // 命令接口
	hooksMixin  // 钩子接口
}

// 下一步待处理的地方
func (c *Client) init() {
	// c.cmdable = c.Process
	// c.initHooks(hooks{
	// 	dial:       c.baseClient.dial,
	// 	process:    c.baseClient.process,
	// 	pipeline:   c.baseClient.processPipeline,
	// 	txPipeline: c.baseClient.processTxPipeline,
	// })
}
