package command

import (
	"context"
	"time"

	"github.com/lufeijun/go-tool-redis/redis/tool/proto"
)

const Nil = proto.Nil

// 执行 redis 命令的基础接口
type Cmder interface {
	Name() string // returns command name in lower case
	FullName() string
	Args() []interface{}
	String() string
	stringArg(int) string
	firstKeyPos() int8
	SetFirstKeyPos(int8)

	readTimeout() *time.Duration
	readReply(rd *proto.Reader) error // read reply from reader

	SetErr(error)
	Err() error
}

type StatefulCmdable interface {
	Cmdable
	Auth(ctx context.Context, password string) *StatusCmd
	AuthACL(ctx context.Context, username, password string) *StatusCmd
	Select(ctx context.Context, index int) *StatusCmd
	SwapDB(ctx context.Context, index1, index2 int) *StatusCmd
	ClientSetName(ctx context.Context, name string) *BoolCmd
	Hello(ctx context.Context, ver int, username, password, clientName string) *MapStringInterfaceCmd
}

// var (
// 	_ Cmdable = (*Client)(nil)
// 	_ Cmdable = (*Tx)(nil)
// 	_ Cmdable = (*Ring)(nil)
// 	_ Cmdable = (*ClusterClient)(nil)
// )

// 所有可以执行的命令
type Cmdable interface {
	Pipeline() Pipeliner
	Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)

	TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)
	TxPipeline() Pipeliner

	Get(ctx context.Context, key string) *StringCmd
}

type Pipeliner interface {
	StatefulCmdable

	// Len is to obtain the number of commands in the pipeline that have not yet been executed.
	Len() int

	// Do is an API for executing any command.
	// If a certain Redis command is not yet supported, you can use Do to execute it.
	Do(ctx context.Context, args ...interface{}) *Cmd

	// Process is to put the commands to be executed into the pipeline buffer.
	Process(ctx context.Context, cmd Cmder) error

	// Discard is to discard all commands in the cache that have not yet been executed.
	Discard()

	// Exec is to send all the commands buffered in the pipeline to the redis-server.
	Exec(ctx context.Context) ([]Cmder, error)
}
