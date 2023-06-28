package pool

import (
	"bufio"
	"net"
	"time"

	"github.com/lufeijun/go-tool-redis/redis/tool/proto"
)

func NewConn(netConn net.Conn) *Conn {
	cn := &Conn{
		netConn:   netConn,
		createdAt: time.Now(),
	}
	cn.rd = proto.NewReader(netConn)
	cn.bw = bufio.NewWriter(netConn)
	cn.wr = proto.NewWriter(cn.bw)
	cn.SetUsedAt(time.Now())
	return cn
}

func NewSingleConnPool(pool Pooler, cn *Conn) *SingleConnPool {
	return &SingleConnPool{
		pool: pool,
		cn:   cn,
	}
}
