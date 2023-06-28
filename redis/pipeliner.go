package redis

import (
	"context"
	"errors"

	"github.com/lufeijun/go-tool-redis/redis/command"
)

type pipelineExecer func(context.Context, []command.Cmder) error

type Pipeliner interface {
	command.StatefulCmdable

	// Len is to obtain the number of commands in the pipeline that have not yet been executed.
	Len() int

	// Do is an API for executing any command.
	// If a certain Redis command is not yet supported, you can use Do to execute it.
	Do(ctx context.Context, args ...interface{}) *command.Cmd

	// Process is to put the commands to be executed into the pipeline buffer.
	Process(ctx context.Context, cmd command.Cmder) error

	// Discard is to discard all commands in the cache that have not yet been executed.
	Discard()

	// Exec is to send all the commands buffered in the pipeline to the redis-server.
	Exec(ctx context.Context) ([]command.Cmder, error)
}

var _ Pipeliner = (*Pipeline)(nil)

type Pipeline struct {
	cmdable
	statefulCmdable

	exec pipelineExecer
	cmds []command.Cmder
}

func (c *Pipeline) init() {
	c.cmdable = c.Process
	c.statefulCmdable = c.Process
}

// 命令个数
func (c *Pipeline) Len() int {
	return len(c.cmds)
}

// Do queues the custom command for later execution.
func (c *Pipeline) Do(ctx context.Context, args ...interface{}) *command.Cmd {
	cmd := command.NewCmd(ctx, args...)
	if len(args) == 0 {
		cmd.SetErr(errors.New("redis: please enter the command to be executed"))
		return cmd
	}
	_ = c.Process(ctx, cmd)
	return cmd
}

// Process queues the cmd for later execution.
func (c *Pipeline) Process(ctx context.Context, cmd command.Cmder) error {
	c.cmds = append(c.cmds, cmd)
	return nil
}

// Discard resets the pipeline and discards queued commands.
func (c *Pipeline) Discard() {
	c.cmds = c.cmds[:0]
}

// Exec executes all previously queued commands using one
// client-server roundtrip.
//
// Exec always returns list of commands and error of the first failed
// command if any.
func (c *Pipeline) Exec(ctx context.Context) ([]command.Cmder, error) {
	if len(c.cmds) == 0 {
		return nil, nil
	}

	cmds := c.cmds
	c.cmds = nil

	return cmds, c.exec(ctx, cmds)
}

func (c *Pipeline) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]command.Cmder, error) {
	if err := fn(c); err != nil {
		return nil, err
	}
	return c.Exec(ctx)
}

func (c *Pipeline) Pipeline() Pipeliner {
	return c
}

func (c *Pipeline) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]command.Cmder, error) {
	return c.Pipelined(ctx, fn)
}

func (c *Pipeline) TxPipeline() Pipeliner {
	return c
}
