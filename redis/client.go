package redis

import "context"

// 可执行命令函数
// type cmdable func(ctx context.Context, cmd command.Cmder) error
// type statefulCmdable func(ctx context.Context, cmd command.Cmder) error

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

func (c *Client) Watch(ctx context.Context, fn func(*Tx) error, keys ...string) error {
	// tx := c.newTx()
	// defer tx.Close(ctx)
	// if len(keys) > 0 {
	// 	if err := tx.Watch(ctx, keys...).Err(); err != nil {
	// 		return err
	// 	}
	// }
	// return fn(tx)
	return nil
}
