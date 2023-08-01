package credential

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/lufeijun/go-tool-wechat/wechat/cache"
	"github.com/lufeijun/go-tool-wechat/wechat/util"
)

const (
	// accessTokenURL 获取access_token的接口
	accessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	// stableAccessTokenURL 获取稳定版access_token的接口
	stableAccessTokenURL = "https://api.weixin.qq.com/cgi-bin/stable_token"
	// workAccessTokenURL 企业微信获取access_token的接口
	workAccessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	// CacheKeyOfficialAccountPrefix 微信公众号cache key前缀
	CacheKeyOfficialAccountPrefix = "lufeijun_gowechat_officialaccount_"
	// CacheKeyMiniProgramPrefix 小程序cache key前缀
	CacheKeyMiniProgramPrefix = "lufeijun_gowechat_miniprogram_"
	// CacheKeyWorkPrefix 企业微信cache key前缀
	CacheKeyWorkPrefix = "lufeijun_gowechat_work_"
)

// DefaultAccessToken 默认AccessToken 获取
type DefaultAccessToken struct {
	appID           string
	appSecret       string
	cacheKeyPrefix  string
	cache           cache.Cache
	accessTokenLock *sync.Mutex
}

// NewDefaultAccessToken new DefaultAccessToken
func NewDefaultAccessToken(appID, appSecret, cacheKeyPrefix string, cache cache.Cache) AccessTokenContextHandle {
	if cache == nil {
		panic("cache is ineed")
	}
	return &DefaultAccessToken{
		appID:           appID,
		appSecret:       appSecret,
		cache:           cache,
		cacheKeyPrefix:  cacheKeyPrefix,
		accessTokenLock: new(sync.Mutex),
	}
}

// ResAccessToken struct
type ResAccessToken struct {
	util.CommonError

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// GetAccessToken 获取access_token,先从cache中获取，没有则从服务端获取
func (ak *DefaultAccessToken) GetAccessToken() (accessToken string, err error) {
	return ak.GetAccessTokenContext(context.Background())
}

// GetAccessTokenContext 获取access_token,先从cache中获取，没有则从服务端获取
func (ak *DefaultAccessToken) GetAccessTokenContext(ctx context.Context) (accessToken string, err error) {
	// 先从cache中取
	// gowechat_officialaccount__access_token_wx82499d23a58480fa = access_token
	accessTokenCacheKey := fmt.Sprintf("%s_access_token_%s", ak.cacheKeyPrefix, ak.appID)
	if val := ak.cache.Get(accessTokenCacheKey); val != nil {
		return val.(string), nil
	}

	// 加上lock，是为了防止在并发获取token时，cache刚好失效，导致从微信服务器上获取到不同token
	ak.accessTokenLock.Lock()
	defer ak.accessTokenLock.Unlock()

	// 双检，防止重复从微信服务器获取
	if val := ak.cache.Get(accessTokenCacheKey); val != nil {
		return val.(string), nil
	}

	// cache失效，从微信服务器获取
	var resAccessToken ResAccessToken
	if resAccessToken, err = GetTokenFromServerContext(ctx, fmt.Sprintf(accessTokenURL, ak.appID, ak.appSecret)); err != nil {
		return
	}

	if err = ak.cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(resAccessToken.ExpiresIn-1500)*time.Second); err != nil {
		return
	}
	accessToken = resAccessToken.AccessToken
	return
}

// GetTokenFromServerContext 强制从微信服务器获取token
func GetTokenFromServerContext(ctx context.Context, url string) (resAccessToken ResAccessToken, err error) {
	var body []byte
	body, err = util.HTTPGetContext(ctx, url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrCode != 0 {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}
	return
}
