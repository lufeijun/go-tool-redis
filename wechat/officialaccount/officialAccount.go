package officialaccount

import (
	"context"

	"github.com/lufeijun/go-tool-wechat/wcontext"
	"github.com/lufeijun/go-tool-wechat/wechat/credential"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/offConfig"
)

// OfficialAccount 微信公众号相关API
type OfficialAccount struct {
	ctx *wcontext.Context
	// basic        *basic.Basic
	// menu         *menu.Menu
	// oauth        *oauth.Oauth
	// material     *material.Material
	// draft        *draft.Draft
	// freepublish  *freepublish.FreePublish
	// js           *js.Js
	// user         *user.User
	// templateMsg  *message.Template
	// managerMsg   *message.Manager
	// device       *device.Device
	// broadcast    *broadcast.Broadcast
	// datacube     *datacube.DataCube
	// ocr          *ocr.OCR
	// subscribeMsg *message.Subscribe
}

// NewOfficialAccount 实例化公众号API
func NewOfficialAccount(cfg *offConfig.Config) *OfficialAccount {
	defaultAkHandle := credential.NewDefaultAccessToken(cfg.AppID, cfg.AppSecret, credential.CacheKeyOfficialAccountPrefix, cfg.Cache)
	ctx := &wcontext.Context{
		Config:            cfg,
		AccessTokenHandle: defaultAkHandle,
	}
	return &OfficialAccount{ctx: ctx}
}

// GetAccessToken 获取access_token
func (officialAccount *OfficialAccount) GetAccessToken() (string, error) {
	return officialAccount.ctx.GetAccessToken()
}

// SetAccessTokenHandle 自定义access_token获取方式
func (officialAccount *OfficialAccount) SetAccessTokenHandle(accessTokenHandle credential.AccessTokenHandle) {
	officialAccount.ctx.AccessTokenHandle = accessTokenHandle
}

// GetAccessTokenContext 获取access_token
func (officialAccount *OfficialAccount) GetAccessTokenContext(ctx context.Context) (string, error) {
	if c, ok := officialAccount.ctx.AccessTokenHandle.(credential.AccessTokenContextHandle); ok {
		return c.GetAccessTokenContext(ctx)
	}
	return officialAccount.ctx.GetAccessToken()
}
