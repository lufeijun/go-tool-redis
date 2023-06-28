package pool

import (
	"bufio"
	"net"
	"sync/atomic"
	"time"

	"github.com/lufeijun/go-tool-redis/redis/tool/proto"
)

type Conn struct {
	usedAt  int64 // atomic
	netConn net.Conn

	rd *proto.Reader
	bw *bufio.Writer
	wr *proto.Writer

	Inited    bool
	pooled    bool
	createdAt time.Time
}

// 链接使用时间
func (cn *Conn) SetUsedAt(tm time.Time) {
	atomic.StoreInt64(&cn.usedAt, tm.Unix())
}
func (cn *Conn) UsedAt() time.Time {
	unix := atomic.LoadInt64(&cn.usedAt)
	return time.Unix(unix, 0)
}

// 关闭连接
func (cn *Conn) Close() error {
	return cn.netConn.Close()
}
