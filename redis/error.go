package redis

import (
	"context"
	"io"
	"net"
	"strings"

	"github.com/lufeijun/go-tool-redis/redis/tool/proto"
)

func isRedisError(err error) bool {
	_, ok := err.(proto.RedisError)
	return ok
}

// 判断是否为错误的连接
func isBadConn(err error, allowTimeout bool, addr string) bool {
	switch err {
	case nil:
		return false
	case context.Canceled, context.DeadlineExceeded:
		return true
	}

	if isRedisError(err) {
		switch {
		case isReadOnlyError(err):
			// Close connections in read only state in case domain addr is used
			// and domain resolves to a different Redis Server. See #790.
			return true
		case isMovedSameConnAddr(err, addr):
			// Close connections when we are asked to move to the same addr
			// of the connection. Force a DNS resolution when all connections
			// of the pool are recycled
			return true
		default:
			return false
		}
	}

	if allowTimeout {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return false
		}
	}

	return true
}

func isReadOnlyError(err error) bool {
	return strings.HasPrefix(err.Error(), "READONLY ")
}

func isMovedSameConnAddr(err error, addr string) bool {
	redisError := err.Error()
	if !strings.HasPrefix(redisError, "MOVED ") {
		return false
	}
	return strings.HasSuffix(redisError, " "+addr)
}

// 是否应该重试
func shouldRetry(err error, retryTimeout bool) bool {
	switch err {
	case io.EOF, io.ErrUnexpectedEOF:
		return true
	case nil, context.Canceled, context.DeadlineExceeded:
		return false
	}

	if v, ok := err.(timeoutError); ok {
		if v.Timeout() {
			return retryTimeout
		}
		return true
	}

	s := err.Error()
	if s == "ERR max number of clients reached" {
		return true
	}
	if strings.HasPrefix(s, "LOADING ") {
		return true
	}
	if strings.HasPrefix(s, "READONLY ") {
		return true
	}
	if strings.HasPrefix(s, "CLUSTERDOWN ") {
		return true
	}
	if strings.HasPrefix(s, "TRYAGAIN ") {
		return true
	}

	return false
}
