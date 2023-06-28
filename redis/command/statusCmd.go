package command

import (
	"context"
	"time"
)

type StatusCmd struct {
	baseCmd

	val string
}

type baseCmd struct {
	ctx    context.Context
	args   []interface{}
	err    error
	keyPos int8

	_readTimeout *time.Duration
}

type Cmd struct {
	baseCmd

	val interface{}
}

func (cmd *baseCmd) SetErr(e error) {
	cmd.err = e
}
