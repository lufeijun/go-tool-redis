package officialaccount

import (
	"context"

	"github.com/lufeijun/go-tool-wechat/wcontext"
	"github.com/lufeijun/go-tool-wechat/wechat/credential"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/basic"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/draft"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/freepublish"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/js"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/material"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/menu"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/oauth"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/offConfig"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/user"
)

// OfficialAccount 微信公众号相关API
type OfficialAccount struct {
	ctx         *wcontext.Context
	basic       *basic.Basic
	menu        *menu.Menu
	oauth       *oauth.Oauth
	material    *material.Material
	draft       *draft.Draft
	freepublish *freepublish.FreePublish
	js          *js.Js
	user        *user.User
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

// GetBasic qr/url 相关配置
func (officialAccount *OfficialAccount) GetBasic() *basic.Basic {
	if officialAccount.basic == nil {
		officialAccount.basic = basic.NewBasic(officialAccount.ctx)
	}
	return officialAccount.basic
}

// GetMenu 菜单管理接口
func (officialAccount *OfficialAccount) GetMenu() *menu.Menu {
	if officialAccount.menu == nil {
		officialAccount.menu = menu.NewMenu(officialAccount.ctx)
	}
	return officialAccount.menu
}

// GetOauth oauth2网页授权
func (officialAccount *OfficialAccount) GetOauth() *oauth.Oauth {
	if officialAccount.oauth == nil {
		officialAccount.oauth = oauth.NewOauth(officialAccount.ctx)
	}
	return officialAccount.oauth
}

// GetMaterial 素材管理
func (officialAccount *OfficialAccount) GetMaterial() *material.Material {
	if officialAccount.material == nil {
		officialAccount.material = material.NewMaterial(officialAccount.ctx)
	}
	return officialAccount.material
}

// GetDraft 草稿箱
func (officialAccount *OfficialAccount) GetDraft() *draft.Draft {
	if officialAccount.draft == nil {
		officialAccount.draft = draft.NewDraft(officialAccount.ctx)
	}
	return officialAccount.draft
}

// GetFreePublish 发布能力
func (officialAccount *OfficialAccount) GetFreePublish() *freepublish.FreePublish {
	if officialAccount.freepublish == nil {
		officialAccount.freepublish = freepublish.NewFreePublish(officialAccount.ctx)
	}
	return officialAccount.freepublish
}

// GetJs js-sdk配置
func (officialAccount *OfficialAccount) GetJs() *js.Js {
	if officialAccount.js == nil {
		officialAccount.js = js.NewJs(officialAccount.ctx)
	}
	return officialAccount.js
}

// GetUser 用户管理接口
func (officialAccount *OfficialAccount) GetUser() *user.User {
	if officialAccount.user == nil {
		officialAccount.user = user.NewUser(officialAccount.ctx)
	}
	return officialAccount.user
}
