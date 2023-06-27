package lufeijun_test

import (
	"testing"
)

func TestRedis(t *testing.T) {

	value := "test"
	if value != "test" {
		t.Error("value failed")
	}
}
