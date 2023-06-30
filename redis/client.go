package redis

import (
	"context"
	"time"

	"github.com/lufeijun/go-tool-redis/redis/pool"
)

// 可执行命令函数

type Client struct {
	*baseClient // 基础客户端字段
	cmdable     // 命令接口
	hooksMixin  // 钩子接口
}

// 下一步待处理的地方
func (c *Client) init() {
	c.cmdable = c.Process
	c.initHooks(hooks{
		dial:       c.baseClient.dial,
		process:    c.baseClient.process,
		pipeline:   c.baseClient.processPipeline,
		txPipeline: c.baseClient.processTxPipeline,
	})
}

// 设置超时时间
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	clone := *c
	clone.baseClient = c.baseClient.withTimeout(timeout)
	clone.init()
	return &clone
}

func (c *Client) Conn() *Conn {
	return newConn(c.opt, pool.NewStickyConnPool(c.connPool))
}

// Do create a Cmd from the args and processes the cmd.
func (c *Client) Do(ctx context.Context, args ...interface{}) *Cmd {
	cmd := NewCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
}

// 命令执行函数
func (c *Client) Process(ctx context.Context, cmd Cmder) error {
	err := c.processHook(ctx, cmd)
	cmd.SetErr(err)
	return err
}

func (c *Client) Options() *Options {
	return c.opt
}

type PoolStats pool.Stats

// 不太明白为什么这么折腾一趟
func (c *Client) PoolStats() *PoolStats {
	stats := c.connPool.Stats()
	return (*PoolStats)(stats)
}

func (c *Client) PoolStats2() *pool.Stats {
	stats := c.connPool.Stats()
	return stats
}

// 管道操作
func (c *Client) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return c.Pipeline().Pipelined(ctx, fn)
}
func (c *Client) Pipeline() Pipeliner {
	pipe := Pipeline{
		exec: pipelineExecer(c.processPipelineHook),
	}
	pipe.init()
	return &pipe
}

// 事务
func (c *Client) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return c.TxPipeline().Pipelined(ctx, fn)
}

// TxPipeline acts like Pipeline, but wraps queued commands with MULTI/EXEC.
func (c *Client) TxPipeline() Pipeliner {
	pipe := Pipeline{
		exec: func(ctx context.Context, cmds []Cmder) error {
			cmds = wrapMultiExec(ctx, cmds)
			return c.processTxPipelineHook(ctx, cmds)
		},
	}
	pipe.init()
	return &pipe
}

// ================ 信息订阅 ========================
// func (c *Client) pubSub() *PubSub {
// 	pubsub := &PubSub{
// 		opt: c.opt,

// 		newConn: func(ctx context.Context, channels []string) (*pool.Conn, error) {
// 			return c.newConn(ctx)
// 		},
// 		closeConn: c.connPool.CloseConn,
// 	}
// 	pubsub.init()
// 	return pubsub
// }

// func (c *Client) Subscribe(ctx context.Context, channels ...string) *PubSub {
// 	pubsub := c.pubSub()
// 	if len(channels) > 0 {
// 		_ = pubsub.Subscribe(ctx, channels...)
// 	}
// 	return pubsub
// }

// PSubscribe subscribes the client to the given patterns.
// Patterns can be omitted to create empty subscription.
// func (c *Client) PSubscribe(ctx context.Context, channels ...string) *PubSub {
// 	pubsub := c.pubSub()
// 	if len(channels) > 0 {
// 		_ = pubsub.PSubscribe(ctx, channels...)
// 	}
// 	return pubsub
// }

// // SSubscribe Subscribes the client to the specified shard channels.
// // Channels can be omitted to create empty subscription.
// func (c *Client) SSubscribe(ctx context.Context, channels ...string) *PubSub {
// 	pubsub := c.pubSub()
// 	if len(channels) > 0 {
// 		_ = pubsub.SSubscribe(ctx, channels...)
// 	}
// 	return pubsub
// }

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
