package command

import (
	"context"

	"github.com/lufeijun/go-tool-redis/redis/tool/proto"
)

type MapStringInterfaceCmd struct {
	baseCmd

	val map[string]interface{}
}

func NewMapStringInterfaceCmd(ctx context.Context, args ...interface{}) *MapStringInterfaceCmd {
	return &MapStringInterfaceCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

var _ Cmder = (*MapStringInterfaceCmd)(nil)

func (cmd *MapStringInterfaceCmd) SetVal(val map[string]interface{}) {
	cmd.val = val
}

func (cmd *MapStringInterfaceCmd) Val() map[string]interface{} {
	return cmd.val
}

func (cmd *MapStringInterfaceCmd) Result() (map[string]interface{}, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *MapStringInterfaceCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *MapStringInterfaceCmd) readReply(rd *proto.Reader) error {
	n, err := rd.ReadMapLen()
	if err != nil {
		return err
	}

	cmd.val = make(map[string]interface{}, n)
	for i := 0; i < n; i++ {
		k, err := rd.ReadString()
		if err != nil {
			return err
		}
		v, err := rd.ReadReply()
		if err != nil {
			if err == Nil {
				cmd.val[k] = Nil
				continue
			}
			if err, ok := err.(proto.RedisError); ok {
				cmd.val[k] = err
				continue
			}
			return err
		}
		cmd.val[k] = v
	}
	return nil
}
