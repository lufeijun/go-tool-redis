package redis

import "context"

type statefulCmdable func(ctx context.Context, cmd Cmder) error

// Hello Set the resp protocol used.
func (c statefulCmdable) Hello(ctx context.Context, ver int, username, password, clientName string) *MapStringInterfaceCmd {
	args := make([]interface{}, 0, 7)
	args = append(args, "hello", ver)
	if password != "" {
		if username != "" {
			args = append(args, "auth", username, password)
		} else {
			args = append(args, "auth", "default", password)
		}
	}
	if clientName != "" {
		args = append(args, "setname", clientName)
	}
	cmd := NewMapStringInterfaceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c statefulCmdable) ClientSetName(ctx context.Context, name string) *BoolCmd {
	cmd := NewBoolCmd(ctx, "client", "setname", name)
	_ = c(ctx, cmd)
	return cmd
}

func (c statefulCmdable) Auth(ctx context.Context, password string) *StatusCmd {
	cmd := NewStatusCmd(ctx, "auth", password)
	_ = c(ctx, cmd)
	return cmd
}

func (c statefulCmdable) AuthACL(ctx context.Context, username, password string) *StatusCmd {
	cmd := NewStatusCmd(ctx, "auth", username, password)
	_ = c(ctx, cmd)
	return cmd
}

func (c statefulCmdable) Select(ctx context.Context, index int) *StatusCmd {
	cmd := NewStatusCmd(ctx, "select", index)
	_ = c(ctx, cmd)
	return cmd
}

func (c statefulCmdable) SwapDB(ctx context.Context, index1, index2 int) *StatusCmd {
	cmd := NewStatusCmd(ctx, "swapdb", index1, index2)
	_ = c(ctx, cmd)
	return cmd
}
