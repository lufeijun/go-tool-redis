package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/lufeijun/go-tool-redis/redis/tool/proto"
	"github.com/lufeijun/go-tool-redis/redis/tool/util"
)

type StringCmd struct {
	baseCmd

	val string
}

var _ Cmder = (*StringCmd)(nil)

func NewStringCmd(ctx context.Context, args ...interface{}) *StringCmd {
	return &StringCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

// 设置值
func (cmd *StringCmd) SetVal(val string) {
	cmd.val = val
}

// 获取值
func (cmd *StringCmd) Val() string {
	return cmd.val
}

// 获取结果 string 类型
func (cmd *StringCmd) Result() (string, error) {
	return cmd.Val(), cmd.err
}

// 获取结果 []byte 类型
func (cmd *StringCmd) Bytes() ([]byte, error) {
	return util.StringToBytes(cmd.val), cmd.err
}

// 获取结果 bool 类型
func (cmd *StringCmd) Bool() (bool, error) {
	if cmd.err != nil {
		return false, cmd.err
	}
	return strconv.ParseBool(cmd.val)
}

// 获取结果 int 类型
func (cmd *StringCmd) Int() (int, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.Atoi(cmd.Val())
}

func (cmd *StringCmd) Int64() (int64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.ParseInt(cmd.Val(), 10, 64)
}

func (cmd *StringCmd) Uint64() (uint64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.ParseUint(cmd.Val(), 10, 64)
}

func (cmd *StringCmd) Float32() (float32, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	f, err := strconv.ParseFloat(cmd.Val(), 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

func (cmd *StringCmd) Float64() (float64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.ParseFloat(cmd.Val(), 64)
}

func (cmd *StringCmd) Time() (time.Time, error) {
	if cmd.err != nil {
		return time.Time{}, cmd.err
	}
	return time.Parse(time.RFC3339Nano, cmd.Val())
}

func (cmd *StringCmd) Scan(val interface{}) error {
	if cmd.err != nil {
		return cmd.err
	}
	return proto.Scan([]byte(cmd.val), val)
}

func (cmd *StringCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *StringCmd) readReply(rd *proto.Reader) (err error) {
	cmd.val, err = rd.ReadString()
	return err
}
