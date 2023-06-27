package pool

import (
	"net"
	"time"
)

type Conn struct {
	usedAt  int64 // atomic
	netConn net.Conn

	// rd *proto.Reader
	// bw *bufio.Writer
	// wr *proto.Writer

	Inited    bool
	pooled    bool
	createdAt time.Time
}
