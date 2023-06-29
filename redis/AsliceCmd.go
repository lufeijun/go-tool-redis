package redis

import (
	"context"

	"github.com/lufeijun/go-tool-redis/redis/tool/hscan"
	"github.com/lufeijun/go-tool-redis/redis/tool/proto"
)

type SliceCmd struct {
	baseCmd

	val []interface{}
}

var _ Cmder = (*SliceCmd)(nil)

func NewSliceCmd(ctx context.Context, args ...interface{}) *SliceCmd {
	return &SliceCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func (cmd *SliceCmd) SetVal(val []interface{}) {
	cmd.val = val
}

func (cmd *SliceCmd) Val() []interface{} {
	return cmd.val
}

func (cmd *SliceCmd) Result() ([]interface{}, error) {
	return cmd.val, cmd.err
}

func (cmd *SliceCmd) String() string {
	return cmdString(cmd, cmd.val)
}

// Scan scans the results from the map into a destination struct. The map keys
// are matched in the Redis struct fields by the `redis:"field"` tag.
func (cmd *SliceCmd) Scan(dst interface{}) error {
	if cmd.err != nil {
		return cmd.err
	}

	// Pass the list of keys and values.
	// Skip the first two args for: HMGET key
	var args []interface{}
	if cmd.args[0] == "hmget" {
		args = cmd.args[2:]
	} else {
		// Otherwise, it's: MGET field field ...
		args = cmd.args[1:]
	}

	return hscan.Scan(dst, args, cmd.val)
}

func (cmd *SliceCmd) readReply(rd *proto.Reader) (err error) {
	cmd.val, err = rd.ReadSlice()
	return err
}
