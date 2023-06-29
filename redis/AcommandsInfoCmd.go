package redis

import (
	"context"
	"fmt"

	"github.com/lufeijun/go-tool-redis/redis/tool/proto"
)

type CommandInfo struct {
	Name        string
	Arity       int8
	Flags       []string
	ACLFlags    []string
	FirstKeyPos int8
	LastKeyPos  int8
	StepCount   int8
	ReadOnly    bool
}

type CommandsInfoCmd struct {
	baseCmd

	val map[string]*CommandInfo
}

var _ Cmder = (*CommandsInfoCmd)(nil)

func NewCommandsInfoCmd(ctx context.Context, args ...interface{}) *CommandsInfoCmd {
	return &CommandsInfoCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func (cmd *CommandsInfoCmd) SetVal(val map[string]*CommandInfo) {
	cmd.val = val
}

func (cmd *CommandsInfoCmd) Val() map[string]*CommandInfo {
	return cmd.val
}

func (cmd *CommandsInfoCmd) Result() (map[string]*CommandInfo, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *CommandsInfoCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *CommandsInfoCmd) readReply(rd *proto.Reader) error {
	const numArgRedis5 = 6
	const numArgRedis6 = 7
	const numArgRedis7 = 10

	n, err := rd.ReadArrayLen()
	if err != nil {
		return err
	}
	cmd.val = make(map[string]*CommandInfo, n)

	for i := 0; i < n; i++ {
		nn, err := rd.ReadArrayLen()
		if err != nil {
			return err
		}

		switch nn {
		case numArgRedis5, numArgRedis6, numArgRedis7:
			// ok
		default:
			return fmt.Errorf("redis: got %d elements in COMMAND reply, wanted 6/7/10", nn)
		}

		cmdInfo := &CommandInfo{}
		if cmdInfo.Name, err = rd.ReadString(); err != nil {
			return err
		}

		arity, err := rd.ReadInt()
		if err != nil {
			return err
		}
		cmdInfo.Arity = int8(arity)

		flagLen, err := rd.ReadArrayLen()
		if err != nil {
			return err
		}
		cmdInfo.Flags = make([]string, flagLen)
		for f := 0; f < len(cmdInfo.Flags); f++ {
			switch s, err := rd.ReadString(); {
			case err == Nil:
				cmdInfo.Flags[f] = ""
			case err != nil:
				return err
			default:
				if !cmdInfo.ReadOnly && s == "readonly" {
					cmdInfo.ReadOnly = true
				}
				cmdInfo.Flags[f] = s
			}
		}

		firstKeyPos, err := rd.ReadInt()
		if err != nil {
			return err
		}
		cmdInfo.FirstKeyPos = int8(firstKeyPos)

		lastKeyPos, err := rd.ReadInt()
		if err != nil {
			return err
		}
		cmdInfo.LastKeyPos = int8(lastKeyPos)

		stepCount, err := rd.ReadInt()
		if err != nil {
			return err
		}
		cmdInfo.StepCount = int8(stepCount)

		if nn >= numArgRedis6 {
			aclFlagLen, err := rd.ReadArrayLen()
			if err != nil {
				return err
			}
			cmdInfo.ACLFlags = make([]string, aclFlagLen)
			for f := 0; f < len(cmdInfo.ACLFlags); f++ {
				switch s, err := rd.ReadString(); {
				case err == Nil:
					cmdInfo.ACLFlags[f] = ""
				case err != nil:
					return err
				default:
					cmdInfo.ACLFlags[f] = s
				}
			}
		}

		if nn >= numArgRedis7 {
			if err := rd.DiscardNext(); err != nil {
				return err
			}
			if err := rd.DiscardNext(); err != nil {
				return err
			}
			if err := rd.DiscardNext(); err != nil {
				return err
			}
		}

		cmd.val[cmdInfo.Name] = cmdInfo
	}

	return nil
}
