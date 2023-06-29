package util

import (
	"context"
	"math/rand"
	"time"
)

func Sleep(ctx context.Context, dur time.Duration) error {
	t := time.NewTimer(dur)
	defer t.Stop()

	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func ToLower(s string) string {
	if isLower(s) {
		return s
	}

	b := make([]byte, len(s))
	for i := range b {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return BytesToString(b)
}

func isLower(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			return false
		}
	}
	return true
}

func RetryBackoff(retry int, minBackoff, maxBackoff time.Duration) time.Duration {
	if retry < 0 {
		panic("not reached")
	}
	if minBackoff == 0 {
		return 0
	}

	d := minBackoff << uint(retry)
	if d < minBackoff {
		return maxBackoff
	}

	d = minBackoff + time.Duration(rand.Int63n(int64(d)))

	if d > maxBackoff || d < minBackoff {
		d = maxBackoff
	}

	return d
}
