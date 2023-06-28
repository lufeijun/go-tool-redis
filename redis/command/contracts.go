package command

import (
	"time"

	"github.com/lufeijun/go-tool-redis/redis/tool/proto"
)

type Cmder interface {
	Name() string
	FullName() string
	Args() []interface{}
	String() string
	stringArg(int) string
	firstKeyPos() int8
	SetFirstKeyPos(int8)

	readTimeout() *time.Duration
	readReply(rd *proto.Reader) error

	SetErr(error)
	Err() error
}

type StatefulCmdable interface {
	// Cmdable
	// Auth(ctx context.Context, password string) *StatusCmd
	// AuthACL(ctx context.Context, username, password string) *StatusCmd
	// Select(ctx context.Context, index int) *StatusCmd
	// SwapDB(ctx context.Context, index1, index2 int) *StatusCmd
	// ClientSetName(ctx context.Context, name string) *BoolCmd
	// Hello(ctx context.Context, ver int, username, password, clientName string) *MapStringInterfaceCmd
}

type Cmdable interface {
}
