package pool

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrClosed performs any operation on the closed client will return this error.
	ErrClosed = errors.New("redis: client is closed")

	// ErrPoolTimeout timed out waiting to get a connection from the connection pool.
	ErrPoolTimeout = errors.New("redis: connection pool timeout")
)

var timers = sync.Pool{
	New: func() interface{} {
		t := time.NewTimer(time.Hour)
		t.Stop()
		return t
	},
}

type lastDialErrorWrap struct {
	err error
}
