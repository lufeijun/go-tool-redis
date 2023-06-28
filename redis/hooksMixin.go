package redis

import (
	"context"
	"net"

	"github.com/lufeijun/go-tool-redis/redis/command"
)

// 给用户自定义的钩子接口
type Hook interface {
	DialHook(next DialHook) DialHook
	ProcessHook(next ProcessHook) ProcessHook
	ProcessPipelineHook(next ProcessPipelineHook) ProcessPipelineHook
}

type (
	DialHook            func(ctx context.Context, network, addr string) (net.Conn, error)
	ProcessHook         func(ctx context.Context, cmd command.Cmder) error
	ProcessPipelineHook func(ctx context.Context, cmds []command.Cmder) error
)

// hooks 部分
type hooks struct {
	dial       DialHook            // 当创建网络连接时调用的hook
	process    ProcessHook         // 执行命令时调用的hook
	pipeline   ProcessPipelineHook // 执行管道命令时调用的hook
	txPipeline ProcessPipelineHook
}

func (h *hooks) setDefaults() {
	if h.dial == nil {
		h.dial = func(ctx context.Context, network, addr string) (net.Conn, error) { return nil, nil }
	}
	if h.process == nil {
		h.process = func(ctx context.Context, cmd command.Cmder) error { return nil }
	}
	if h.pipeline == nil {
		h.pipeline = func(ctx context.Context, cmds []command.Cmder) error { return nil }
	}
	if h.txPipeline == nil {
		h.txPipeline = func(ctx context.Context, cmds []command.Cmder) error { return nil }
	}
}

// 接口实现
type hooksMixin struct {
	slice   []Hook
	initial hooks // 初始
	current hooks // 当前
}

// 初始化
func (hs *hooksMixin) initHooks(hooks hooks) {
	hs.initial = hooks
	hs.chain()
}

// 链？
func (hs *hooksMixin) chain() {
	// 内部 hooks 初始化
	hs.initial.setDefaults()

	// 重置当前 hooks
	hs.current.dial = hs.initial.dial
	hs.current.process = hs.initial.process
	hs.current.pipeline = hs.initial.pipeline
	hs.current.txPipeline = hs.initial.txPipeline

	// 遍历 slice
	for i := len(hs.slice) - 1; i >= 0; i-- {
		if wrapped := hs.slice[i].DialHook(hs.current.dial); wrapped != nil {
			hs.current.dial = wrapped
		}
		if wrapped := hs.slice[i].ProcessHook(hs.current.process); wrapped != nil {
			hs.current.process = wrapped
		}
		if wrapped := hs.slice[i].ProcessPipelineHook(hs.current.pipeline); wrapped != nil {
			hs.current.pipeline = wrapped
		}
		if wrapped := hs.slice[i].ProcessPipelineHook(hs.current.txPipeline); wrapped != nil {
			hs.current.txPipeline = wrapped
		}
	}
}

// 添加 hook
func (hs *hooksMixin) AddHook(hook Hook) {
	hs.slice = append(hs.slice, hook)
	hs.chain()
}

// 克隆
func (hs *hooksMixin) clone() hooksMixin {
	clone := *hs
	l := len(clone.slice)
	clone.slice = clone.slice[:l:l]
	return clone
}

// redis 命令钩子函数
func (hs *hooksMixin) withProcessHook(ctx context.Context, cmd command.Cmder, hook ProcessHook) error {

	// 获取了 slice 最后一个满足条件的元素
	for i := len(hs.slice) - 1; i >= 0; i-- {
		if wrapped := hs.slice[i].ProcessHook(hook); wrapped != nil {
			hook = wrapped
		}
	}
	return hook(ctx, cmd)
}

// redis 管道命令钩子函数
func (hs *hooksMixin) withProcessPipelineHook(ctx context.Context, cmds []command.Cmder, hook ProcessPipelineHook) error {
	for i := len(hs.slice) - 1; i >= 0; i-- {
		if wrapped := hs.slice[i].ProcessPipelineHook(hook); wrapped != nil {
			hook = wrapped
		}
	}
	return hook(ctx, cmds)
}

// 调用当前的钩子函数
func (hs *hooksMixin) dialHook(ctx context.Context, network, addr string) (net.Conn, error) {
	return hs.current.dial(ctx, network, addr)
}
func (hs *hooksMixin) processHook(ctx context.Context, cmd command.Cmder) error {
	return hs.current.process(ctx, cmd)
}

func (hs *hooksMixin) processPipelineHook(ctx context.Context, cmds []command.Cmder) error {
	return hs.current.pipeline(ctx, cmds)
}

func (hs *hooksMixin) processTxPipelineHook(ctx context.Context, cmds []command.Cmder) error {
	return hs.current.txPipeline(ctx, cmds)
}
