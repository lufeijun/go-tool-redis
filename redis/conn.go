package redis

import (
	"context"

	"github.com/lufeijun/go-tool-redis/redis/command"
)

type Conn struct {
	baseClient
	cmdable
	statefulCmdable
	hooksMixin
}

func (c *Conn) Process(ctx context.Context, cmd command.Cmder) error {
	err := c.processHook(ctx, cmd)
	cmd.SetErr(err)
	return err
}

func (c *Conn) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]command.Cmder, error) {
	return c.Pipeline().Pipelined(ctx, fn)
}

func (c *Conn) Pipeline() Pipeliner {
	pipe := Pipeline{
		exec: c.processPipelineHook,
	}
	pipe.init()
	return &pipe
}
