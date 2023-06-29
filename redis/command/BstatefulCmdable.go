package command

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
