package redis

import (
	"context"

	"github.com/lufeijun/go-tool-redis/redis/tool/proto"
	"github.com/lufeijun/go-tool-redis/redis/tool/util"
)

// bool 类型的命令
type BoolCmd struct {
	baseCmd

	val bool
}

var _ Cmder = (*BoolCmd)(nil)

func NewBoolCmd(ctx context.Context, args ...interface{}) *BoolCmd {
	return &BoolCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func (cmd *BoolCmd) SetVal(val bool) {
	cmd.val = val
}

func (cmd *BoolCmd) Val() bool {
	return cmd.val
}

func (cmd *BoolCmd) Result() (bool, error) {
	return cmd.val, cmd.err
}

func (cmd *BoolCmd) String() string {
	return cmdString(cmd, cmd.val)
}

// 读取结果
func (cmd *BoolCmd) readReply(rd *proto.Reader) (err error) {
	cmd.val, err = rd.ReadBool()

	// 兼容性处理
	// `SET key value NX` returns nil when key already exists. But
	// `SETNX key value` returns bool (0/1). So convert nil to bool.
	if err == Nil {
		cmd.val = false
		err = nil
	}
	return err
}

func cmdString(cmd Cmder, val interface{}) string {
	b := make([]byte, 0, 64)

	for i, arg := range cmd.Args() {
		if i > 0 {
			b = append(b, ' ')
		}
		b = util.AppendArg(b, arg)
	}

	if err := cmd.Err(); err != nil {
		b = append(b, ": "...)
		b = append(b, err.Error()...)
	} else if val != nil {
		b = append(b, ": "...)
		b = util.AppendArg(b, val)
	}

	return util.BytesToString(b)
}
