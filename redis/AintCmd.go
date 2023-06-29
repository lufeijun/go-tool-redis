package redis

import (
	"context"

	"github.com/lufeijun/go-tool-redis/redis/tool/proto"
)

type IntCmd struct {
	baseCmd

	val int64
}

var _ Cmder = (*IntCmd)(nil)

func NewIntCmd(ctx context.Context, args ...interface{}) *IntCmd {
	return &IntCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func (cmd *IntCmd) SetVal(val int64) {
	cmd.val = val
}

func (cmd *IntCmd) Val() int64 {
	return cmd.val
}

func (cmd *IntCmd) Result() (int64, error) {
	return cmd.val, cmd.err
}

func (cmd *IntCmd) Uint64() (uint64, error) {
	return uint64(cmd.val), cmd.err
}

func (cmd *IntCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *IntCmd) readReply(rd *proto.Reader) (err error) {
	cmd.val, err = rd.ReadInt()
	return err
}
