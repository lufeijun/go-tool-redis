package tcb

import "github.com/lufeijun/go-tool-wechat/wechat/miniprogram/wcontext"

// Tcb Tencent Cloud Base
type Tcb struct {
	*wcontext.Context
}

// NewTcb new Tencent Cloud Base
func NewTcb(context *wcontext.Context) *Tcb {
	return &Tcb{
		context,
	}
}
