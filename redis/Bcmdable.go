package redis

import (
	"context"
	"time"
)

type cmdable func(ctx context.Context, cmd Cmder) error

func (c cmdable) Wait(ctx context.Context, numSlaves int, timeout time.Duration) *IntCmd {
	cmd := NewIntCmd(ctx, "wait", numSlaves, int(timeout/time.Millisecond))
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Command(ctx context.Context) *CommandsInfoCmd {
	cmd := NewCommandsInfoCmd(ctx, "command")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Get(ctx context.Context, key string) *StringCmd {
	cmd := NewStringCmd(ctx, "get", key)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ReadOnly(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd(ctx, "readonly")
	_ = c(ctx, cmd)
	return cmd
}
