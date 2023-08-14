package work

import (
	"github.com/lufeijun/go-tool-wechat/wechat/credential"
	"github.com/lufeijun/go-tool-wechat/wechat/work/wcontext"
	"github.com/lufeijun/go-tool-wechat/wechat/work/workconfig"
)

// Work 企业微信
type Work struct {
	ctx *wcontext.Context
}

// NewWork init work
func NewWork(cfg *workconfig.Config) *Work {
	defaultAkHandle := credential.NewWorkAccessToken(cfg.CorpID, cfg.CorpSecret, credential.CacheKeyWorkPrefix, cfg.Cache)
	ctx := &wcontext.Context{
		Config:            cfg,
		AccessTokenHandle: defaultAkHandle,
	}
	return &Work{ctx: ctx}
}

// GetContext get Context
func (wk *Work) GetContext() *wcontext.Context {
	return wk.ctx
}

// // GetOauth get oauth
// func (wk *Work) GetOauth() *oauth.Oauth {
// 	return oauth.NewOauth(wk.ctx)
// }

// // GetMsgAudit get msgAudit
// func (wk *Work) GetMsgAudit() (*msgaudit.Client, error) {
// 	return msgaudit.NewClient(wk.ctx.Config)
// }

// // GetKF get kf
// func (wk *Work) GetKF() (*kf.Client, error) {
// 	return kf.NewClient(wk.ctx.Config)
// }

// // GetExternalContact get external_contact
// func (wk *Work) GetExternalContact() *externalcontact.Client {
// 	return externalcontact.NewClient(wk.ctx)
// }

// // GetAddressList get address_list
// func (wk *Work) GetAddressList() *addresslist.Client {
// 	return addresslist.NewClient(wk.ctx)
// }

// // GetMaterial get material
// func (wk *Work) GetMaterial() *material.Client {
// 	return material.NewClient(wk.ctx)
// }

// // GetRobot get robot
// func (wk *Work) GetRobot() *robot.Client {
// 	return robot.NewClient(wk.ctx)
// }

// // GetMessage 获取发送应用消息接口实例
// func (wk *Work) GetMessage() *message.Client {
// 	return message.NewClient(wk.ctx)
// }

// // GetAppChat 获取应用发送消息到群聊会话接口实例
// func (wk *Work) GetAppChat() *appchat.Client {
// 	return appchat.NewClient(wk.ctx)
// }

// // GetInvoice get invoice
// func (wk *Work) GetInvoice() *invoice.Client {
// 	return invoice.NewClient(wk.ctx)
// }
