package openapi

import (
	"errors"
	"fmt"

	ocContext "github.com/lufeijun/go-tool-wechat/wechat/officialaccount/wcontext"
	"github.com/lufeijun/go-tool-wechat/wechat/util"
)

const (
	clearQuotaURL            = "https://api.weixin.qq.com/cgi-bin/clear_quota"       // 重置API调用次数
	getAPIQuotaURL           = "https://api.weixin.qq.com/cgi-bin/openapi/quota/get" // 查询API调用额度
	getRidInfoURL            = "https://api.weixin.qq.com/cgi-bin/openapi/rid/get"   // 查询rid信息
	clearQuotaByAppSecretURL = "https://api.weixin.qq.com/cgi-bin/clear_quota/v2"    // 使用AppSecret重置 API 调用次数
)

// OpenAPI openApi管理
type OpenAPI struct {
	ctx interface{}
}

// GetAPIQuotaParams 查询API调用额度参数
type GetAPIQuotaParams struct {
	CgiPath string `json:"cgi_path"` // api的请求地址，例如"/cgi-bin/message/custom/send";不要前缀“https://api.weixin.qq.com” ，也不要漏了"/",否则都会76003的报错
}

// APIQuota API调用额度
type APIQuota struct {
	util.CommonError
	Quota struct {
		DailyLimit int64 `json:"daily_limit"` // 当天该账号可调用该接口的次数
		Used       int64 `json:"used"`        // 当天已经调用的次数
		Remain     int64 `json:"remain"`      // 当天剩余调用次数
	} `json:"quota"` // 详情
}

// GetRidInfoParams 查询rid信息参数
type GetRidInfoParams struct {
	Rid string `json:"rid"` // 调用接口报错返回的rid
}

// RidInfo rid信息
type RidInfo struct {
	util.CommonError
	Request struct {
		InvokeTime   int64  `json:"invoke_time"`   // 发起请求的时间戳
		CostInMs     int64  `json:"cost_in_ms"`    // 请求毫秒级耗时
		RequestURL   string `json:"request_url"`   // 请求的URL参数
		RequestBody  string `json:"request_body"`  // post请求的请求参数
		ResponseBody string `json:"response_body"` // 接口请求返回参数
		ClientIP     string `json:"client_ip"`     // 接口请求的客户端ip
	} `json:"request"` // 该rid对应的请求详情
}

// NewOpenAPI 实例化
func NewOpenAPI(ctx interface{}) *OpenAPI {
	return &OpenAPI{ctx: ctx}
}

// ClearQuota 重置API调用次数
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/openApi-mgnt/clearQuota.html
func (o *OpenAPI) ClearQuota() error {
	appID, _, err := o.getAppIDAndSecret()
	if err != nil {
		return err
	}

	var payload = struct {
		AppID string `json:"appid"`
	}{
		AppID: appID,
	}
	res, err := o.doPostRequest(clearQuotaURL, payload)
	if err != nil {
		return err
	}

	return util.DecodeWithCommonError(res, "ClearQuota")
}

// GetAPIQuota 查询API调用额度
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/openApi-mgnt/getApiQuota.html
func (o *OpenAPI) GetAPIQuota(params GetAPIQuotaParams) (quota APIQuota, err error) {
	res, err := o.doPostRequest(getAPIQuotaURL, params)
	if err != nil {
		return
	}

	err = util.DecodeWithError(res, &quota, "GetAPIQuota")
	return
}

// GetRidInfo 查询rid信息
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/openApi-mgnt/getRidInfo.html
func (o *OpenAPI) GetRidInfo(params GetRidInfoParams) (r RidInfo, err error) {
	res, err := o.doPostRequest(getRidInfoURL, params)
	if err != nil {
		return
	}

	err = util.DecodeWithError(res, &r, "GetRidInfo")
	return
}

// ClearQuotaByAppSecret 使用AppSecret重置 API 调用次数
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/openApi-mgnt/clearQuotaByAppSecret.html
func (o *OpenAPI) ClearQuotaByAppSecret() error {
	id, secret, err := o.getAppIDAndSecret()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?appid=%s&appsecret=%s", clearQuotaByAppSecretURL, id, secret)
	res, err := util.HTTPPost(uri, "")
	if err != nil {
		return err
	}

	return util.DecodeWithCommonError(res, "ClearQuotaByAppSecret")
}

// 获取 AppID 和 AppSecret
func (o *OpenAPI) getAppIDAndSecret() (string, string, error) {
	switch o.ctx.(type) {
	// case *mpContext.Context:
	// 	c := o.ctx.(*mpContext.Context)
	// 	return c.AppID, c.AppSecret, nil
	case *ocContext.Context:
		c := o.ctx.(*ocContext.Context)
		return c.AppID, c.AppSecret, nil
	default:
		return "", "", errors.New("invalid context type")
	}
}

// 获取 AccessToken
func (o *OpenAPI) getAccessToken() (string, error) {
	switch o.ctx.(type) {
	// case *mpContext.Context:
	// 	return o.ctx.(*mpContext.Context).GetAccessToken()
	case *ocContext.Context:
		return o.ctx.(*ocContext.Context).GetAccessToken()
	default:
		return "", errors.New("invalid context type")
	}
}

// 创建 POST 请求
func (o *OpenAPI) doPostRequest(uri string, payload interface{}) ([]byte, error) {
	ak, err := o.getAccessToken()
	if err != nil {
		return nil, err
	}

	uri = fmt.Sprintf("%s?access_token=%s", uri, ak)
	return util.PostJSON(uri, payload)
}
