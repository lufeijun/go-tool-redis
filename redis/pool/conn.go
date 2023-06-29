package pool

import (
	"bufio"
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/lufeijun/go-tool-redis/redis/tool/proto"
)

var noDeadline = time.Time{}

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

// 发送命令
func (cn *Conn) WithWriter(
	ctx context.Context, timeout time.Duration, fn func(wr *proto.Writer) error,
) error {
	if timeout >= 0 {
		if err := cn.netConn.SetWriteDeadline(cn.deadline(ctx, timeout)); err != nil {
			return err
		}
	}

	// 重置缓冲区
	if cn.bw.Buffered() > 0 {
		cn.bw.Reset(cn.netConn)
	}

	if err := fn(cn.wr); err != nil {
		return err
	}

	return cn.bw.Flush()
}

// 设置超时时间
func (cn *Conn) deadline(ctx context.Context, timeout time.Duration) time.Time {
	tm := time.Now()
	cn.SetUsedAt(tm)

	if timeout > 0 {
		tm = tm.Add(timeout)
	}

	// 比较两个时间
	if ctx != nil {
		deadline, ok := ctx.Deadline()
		if ok {
			if timeout == 0 {
				return deadline
			}
			if deadline.Before(tm) {
				return deadline
			}
			return tm
		}
	}

	if timeout > 0 {
		return tm
	}

	return noDeadline
}

// 接受返回值
func (cn *Conn) WithReader(
	ctx context.Context, timeout time.Duration, fn func(rd *proto.Reader) error,
) error {
	if timeout >= 0 {
		if err := cn.netConn.SetReadDeadline(cn.deadline(ctx, timeout)); err != nil {
			return err
		}
	}
	return fn(cn.rd)
}
