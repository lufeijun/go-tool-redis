package main_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/lufeijun/go-tool-godoten/env"
)

func TestAaa(t *testing.T) {
	err := env.Load()

	if err != nil {
		panic(err)
	}

	bbb := os.Getenv("bbb")

	fmt.Println(bbb)
}
