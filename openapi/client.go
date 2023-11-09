package openapi

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/lufeijun/go-tool-aliyun/common"
	"github.com/lufeijun/go-tool-aliyun/credential"
	"github.com/lufeijun/go-tool-aliyun/service"
)

type Client struct {
	Endpoint             *string
	RegionId             *string
	Protocol             *string
	Method               *string
	UserAgent            *string
	EndpointRule         *string
	EndpointMap          map[string]*string
	Suffix               *string
	ReadTimeout          *int
	ConnectTimeout       *int
	HttpProxy            *string
	HttpsProxy           *string
	Socks5Proxy          *string
	Socks5NetWork        *string
	NoProxy              *string
	Network              *string
	ProductId            *string
	MaxIdleConns         *int
	EndpointType         *string
	OpenPlatformEndpoint *string
	Credential           credential.Credential
	SignatureVersion     *string
	SignatureAlgorithm   *string
	Headers              map[string]*string
	// Spi                spi.ClientInterface
	GlobalParameters *GlobalParameters
	Key              *string
	Cert             *string
	Ca               *string
}

/**
 * Init client with Config
 * @param config config contains the necessary information to create a client
 */
func NewClient(config *Config) (*Client, error) {
	client := new(Client)
	err := client.Init(config)
	return client, err
}

func (client *Client) Init(config *Config) (_err error) {
	// 参数检测
	if common.BoolValue(common.IsUnset(config)) {
		_err = common.NewSDKError(map[string]interface{}{
			"code":    "ParameterMissing",
			"message": "'config' can not be unset",
		})
		return
	}

	if !common.BoolValue(common.Empty(config.AccessKeyId)) && !common.BoolValue(common.Empty(config.AccessKeySecret)) {
		if !common.BoolValue(common.Empty(config.SecurityToken)) {
			config.Type = common.String("sts")
		} else {
			config.Type = common.String("access_key")
		}

		credentialConfig := &credential.Config{
			AccessKeyId:     config.AccessKeyId,
			Type:            config.Type,
			AccessKeySecret: config.AccessKeySecret,
		}
		credentialConfig.SecurityToken = config.SecurityToken
		client.Credential, _err = credential.NewCredential(credentialConfig)
		if _err != nil {
			return _err
		}

	} else if !common.BoolValue(common.IsUnset(config.Credential)) {
		client.Credential = config.Credential
	}

	// 基础字段赋值
	client.Endpoint = config.Endpoint
	client.EndpointType = config.EndpointType
	client.Network = config.Network
	client.Suffix = config.Suffix
	client.Protocol = config.Protocol
	client.Method = config.Method
	client.RegionId = config.RegionId
	client.UserAgent = config.UserAgent
	client.ReadTimeout = config.ReadTimeout
	client.ConnectTimeout = config.ConnectTimeout
	client.HttpProxy = config.HttpProxy
	client.HttpsProxy = config.HttpsProxy
	client.NoProxy = config.NoProxy
	client.Socks5Proxy = config.Socks5Proxy
	client.Socks5NetWork = config.Socks5NetWork
	client.MaxIdleConns = config.MaxIdleConns
	client.SignatureVersion = config.SignatureVersion
	client.SignatureAlgorithm = config.SignatureAlgorithm
	client.GlobalParameters = config.GlobalParameters
	client.Key = config.Key
	client.Cert = config.Cert
	client.Ca = config.Ca
	return
}

func (client *Client) CallApi(params *Params, request *OpenApiRequest, runtime *service.RuntimeOptions) (_result map[string]interface{}, _err error) {

	if common.BoolValue(common.IsUnset(params)) {
		_err = common.NewSDKError(map[string]interface{}{
			"code":    "ParameterMissing",
			"message": "'params' 缺少请求参数",
		})
		return
	}

	if common.BoolValue(common.IsUnset(client.SignatureAlgorithm)) || !common.BoolValue(common.EqualString(client.SignatureAlgorithm, common.String("v2"))) {
		_result = make(map[string]interface{})
		_body, _err := client.DoRequest(params, request, runtime)
		if _err != nil {
			return _result, _err
		}
		_result = _body
		return _result, _err
	}

	// else if common.BoolValue(common.EqualString(params.Style, common.String("ROA"))) && common.BoolValue(common.EqualString(params.ReqBodyType, common.String("json"))) {
	// 	_result = make(map[string]interface{})
	// 	_body, _err := client.DoROARequest(params.Action, params.Version, params.Protocol, params.Method, params.AuthType, params.Pathname, params.BodyType, request, runtime)
	// 	if _err != nil {
	// 		return _result, _err
	// 	}
	// 	_result = _body
	// 	return _result, _err
	// } else if common.BoolValue(common.EqualString(params.Style, common.String("ROA"))) {
	// 	_result = make(map[string]interface{})
	// 	_body, _err := client.DoROARequestWithForm(params.Action, params.Version, params.Protocol, params.Method, params.AuthType, params.Pathname, params.BodyType, request, runtime)
	// 	if _err != nil {
	// 		return _result, _err
	// 	}
	// 	_result = _body
	// 	return _result, _err
	// } else { 下次整这块
	// 	_result = make(map[string]interface{})
	// 	_body, _err := client.DoRPCRequest(params.Action, params.Version, params.Protocol, params.Method, params.AuthType, params.BodyType, request, runtime)
	// 	if _err != nil {
	// 		return _result, _err
	// 	}
	// 	_result = _body
	// 	return _result, _err
	// }

	return
}

func (client *Client) DoRequest(params *Params, request *OpenApiRequest, runtime *service.RuntimeOptions) (_result map[string]interface{}, _err error) {
	_err = common.Validate(params)
	if _err != nil {
		return _result, _err
	}
	_err = common.Validate(request)
	if _err != nil {
		return _result, _err
	}
	_err = common.Validate(runtime)
	if _err != nil {
		return _result, _err
	}

	_runtime := map[string]interface{}{
		"timeouted":      "retry",
		"key":            common.StringValue(common.DefaultString(runtime.Key, client.Key)),
		"cert":           common.StringValue(common.DefaultString(runtime.Cert, client.Cert)),
		"ca":             common.StringValue(common.DefaultString(runtime.Ca, client.Ca)),
		"readTimeout":    common.IntValue(common.DefaultNumber(runtime.ReadTimeout, client.ReadTimeout)),
		"connectTimeout": common.IntValue(common.DefaultNumber(runtime.ConnectTimeout, client.ConnectTimeout)),
		"httpProxy":      common.StringValue(common.DefaultString(runtime.HttpProxy, client.HttpProxy)),
		"httpsProxy":     common.StringValue(common.DefaultString(runtime.HttpsProxy, client.HttpsProxy)),
		"noProxy":        common.StringValue(common.DefaultString(runtime.NoProxy, client.NoProxy)),
		"socks5Proxy":    common.StringValue(common.DefaultString(runtime.Socks5Proxy, client.Socks5Proxy)),
		"socks5NetWork":  common.StringValue(common.DefaultString(runtime.Socks5NetWork, client.Socks5NetWork)),
		"maxIdleConns":   common.IntValue(common.DefaultNumber(runtime.MaxIdleConns, client.MaxIdleConns)),
		"retry": map[string]interface{}{
			"retryable":   common.BoolValue(runtime.Autoretry),
			"maxAttempts": common.IntValue(common.DefaultNumber(runtime.MaxAttempts, common.Int(3))),
		},
		"backoff": map[string]interface{}{
			"policy": common.StringValue(common.DefaultString(runtime.BackoffPolicy, common.String("no"))),
			"period": common.IntValue(common.DefaultNumber(runtime.BackoffPeriod, common.Int(1))),
		},
		"ignoreSSL": common.BoolValue(runtime.IgnoreSSL),
	}

	_resp := make(map[string]interface{})

	// 请求重试 start
	for _retryTimes := 0; common.BoolValue(common.AllowRetry(_runtime["retry"], common.Int(_retryTimes))); _retryTimes++ {
		// 重试机制处理
		if _retryTimes > 0 {
			// 已经是重试了，需要休眠等待一会
			_backoffTime := common.GetBackoffTime(_runtime["backoff"], common.Int(_retryTimes))
			if common.IntValue(_backoffTime) > 0 {
				common.Sleep(_backoffTime)
			}
		}

		// 处理请求 start
		_resp, _err = func() (map[string]interface{}, error) {
			request_ := common.NewRequest()
			request_.Protocol = common.DefaultString(client.Protocol, params.Protocol)
			request_.Method = params.Method
			request_.Pathname = params.Pathname
			globalQueries := make(map[string]*string)
			globalHeaders := make(map[string]*string)

			if !common.BoolValue(common.IsUnset(client.GlobalParameters)) {
				globalParams := client.GlobalParameters
				if !common.BoolValue(common.IsUnset(globalParams.Queries)) {
					globalQueries = globalParams.Queries
				}

				if !common.BoolValue(common.IsUnset(globalParams.Headers)) {
					globalHeaders = globalParams.Headers
				}
			}

			request_.Query = common.Merge(globalQueries, request.Query)

			// endpoint is setted in product client
			request_.Headers = common.Merge(map[string]*string{
				"host":                  client.Endpoint,
				"x-acs-version":         params.Version,
				"x-acs-action":          params.Action,
				"user-agent":            service.GetUserAgent(client.UserAgent),
				"x-acs-date":            service.GetTimestamp(),
				"x-acs-signature-nonce": service.GetNonce(),
				"accept":                common.String("application/json"),
			}, globalHeaders, request.Headers)

			if common.BoolValue(common.EqualString(params.Style, common.String("RPC"))) {
				headers, _err := client.GetRpcHeaders()
				if _err != nil {
					return _result, _err
				}

				if !common.BoolValue(common.IsUnset(headers)) {
					request_.Headers = common.Merge(request_.Headers,
						headers)
				}
			}

			signatureAlgorithm := common.DefaultString(client.SignatureAlgorithm, common.String("ACS3-HMAC-SHA256"))
			hashedRequestPayload := service.HexEncode(service.Hash(service.ToBytes(common.String("")), signatureAlgorithm))

			if !common.BoolValue(common.IsUnset(request.Stream)) {
				tmp, _err := service.ReadAsBytes(request.Stream)
				if _err != nil {
					return _result, _err
				}
				hashedRequestPayload = service.HexEncode(service.Hash(tmp, signatureAlgorithm))
				request_.Body = common.ToReader(tmp)
				request_.Headers["content-type"] = common.String("application/octet-stream")
			} else {
				if !common.BoolValue(common.IsUnset(request.Body)) {
					if common.BoolValue(common.EqualString(params.ReqBodyType, common.String("json"))) {
						jsonObj := service.ToJSONString(request.Body)
						hashedRequestPayload = service.HexEncode(service.Hash(service.ToBytes(jsonObj), signatureAlgorithm))
						request_.Body = common.ToReader(jsonObj)
						request_.Headers["content-type"] = common.String("application/json; charset=utf-8")
					} else {
						m, _err := service.AssertAsMap(request.Body)
						if _err != nil {
							return _result, _err
						}

						formObj := service.ToForm(m)
						hashedRequestPayload = service.HexEncode(service.Hash(service.ToBytes(formObj), signatureAlgorithm))
						request_.Body = common.ToReader(formObj)
						request_.Headers["content-type"] = common.String("application/x-www-form-urlencoded")
					}

				}
			}

			//
			request_.Headers["x-acs-content-sha256"] = hashedRequestPayload

			if !common.BoolValue(service.EqualString(params.AuthType, common.String("Anonymous"))) {
				authType, _err := client.GetType()
				if _err != nil {
					return _result, _err
				}

				if common.BoolValue(common.EqualString(authType, common.String("bearer"))) {
					bearerToken, _err := client.GetBearerToken()
					if _err != nil {
						return _result, _err
					}

					request_.Headers["x-acs-bearer-token"] = bearerToken
				} else {
					accessKeyId, _err := client.GetAccessKeyId()
					if _err != nil {
						return _result, _err
					}

					accessKeySecret, _err := client.GetAccessKeySecret()
					if _err != nil {
						return _result, _err
					}

					securityToken, _err := client.GetSecurityToken()
					if _err != nil {
						return _result, _err
					}

					if !common.BoolValue(common.Empty(securityToken)) {
						request_.Headers["x-acs-accesskey-id"] = accessKeyId
						request_.Headers["x-acs-security-token"] = securityToken
					}

					request_.Headers["Authorization"] = service.GetAuthorization(request_, signatureAlgorithm, hashedRequestPayload, accessKeyId, accessKeySecret)
				}

			}

			//
			response_, _err := common.DoRequest(request_, _runtime)
			if _err != nil {
				return _result, _err
			}

			// 错误处理
			if common.BoolValue(service.Is4xx(response_.StatusCode)) || common.BoolValue(service.Is5xx(response_.StatusCode)) {
				err := map[string]interface{}{}
				if !common.BoolValue(common.IsUnset(response_.Headers["content-type"])) && common.BoolValue(common.EqualString(response_.Headers["content-type"], common.String("text/xml;charset=utf-8"))) {
					// _str, _err := service.ReadAsString(response_.Body)
					// if _err != nil {
					// 	return _result, _err
					// }

					// respMap := xml.ParseXml(_str, nil)
					// err, _err = service.AssertAsMap(respMap["Error"])
					// if _err != nil {
					// 	return _result, _err
					// }

				} else {
					_res, _err := service.ReadAsJSON(response_.Body)
					if _err != nil {
						return _result, _err
					}

					err, _err = service.AssertAsMap(_res)
					if _err != nil {
						return _result, _err
					}

				}

				err["statusCode"] = response_.StatusCode
				_err = common.NewSDKError(map[string]interface{}{
					"code":               common.ToString(DefaultAny(err["Code"], err["code"])),
					"message":            "code: " + common.ToString(tea.IntValue(response_.StatusCode)) + ", " + tea.ToString(DefaultAny(err["Message"], err["message"])) + " request id: " + tea.ToString(DefaultAny(err["RequestId"], err["requestId"])),
					"data":               err,
					"description":        tea.ToString(DefaultAny(err["Description"], err["description"])),
					"accessDeniedDetail": DefaultAny(err["AccessDeniedDetail"], err["accessDeniedDetail"]),
				})
				return _result, _err
			}

			// 处理结果
			if common.BoolValue(common.EqualString(params.BodyType, common.String("binary"))) {
				resp := map[string]interface{}{
					"body":       response_.Body,
					"headers":    response_.Headers,
					"statusCode": common.IntValue(response_.StatusCode),
				}
				_result = resp
				return _result, _err
			} else if common.BoolValue(common.EqualString(params.BodyType, common.String("byte"))) {
				byt, _err := service.ReadAsBytes(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = common.Convert(map[string]interface{}{
					"body":       byt,
					"headers":    response_.Headers,
					"statusCode": common.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if common.BoolValue(common.EqualString(params.BodyType, common.String("string"))) {
				str, _err := service.ReadAsString(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = common.Convert(map[string]interface{}{
					"body":       common.StringValue(str),
					"headers":    response_.Headers,
					"statusCode": common.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if common.BoolValue(common.EqualString(params.BodyType, common.String("json"))) {
				obj, _err := service.ReadAsJSON(response_.Body)
				if _err != nil {
					return _result, _err
				}

				res, _err := service.AssertAsMap(obj)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       res,
					"headers":    response_.Headers,
					"statusCode": common.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if common.BoolValue(common.EqualString(params.BodyType, common.String("array"))) {
				arr, _err := service.ReadAsJSON(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = common.Convert(map[string]interface{}{
					"body":       arr,
					"headers":    response_.Headers,
					"statusCode": common.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else {
				_result = make(map[string]interface{})
				_err = common.Convert(map[string]interface{}{
					"headers":    response_.Headers,
					"statusCode": common.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			}

		}()
		// 处理请求 end

		if !common.BoolValue(common.Retryable(_err)) {
			break
		}

	}
	// 请求重试 end

	return _resp, _err
}

/**
 * get RPC header for debug
 */
func (client *Client) GetRpcHeaders() (_result map[string]*string, _err error) {
	headers := client.Headers
	client.Headers = nil
	_result = headers
	return _result, _err
}

/**
 * Get credential type by credential
 * @return credential type e.g. access_key
 */
func (client *Client) GetType() (*string, error) {
	if common.BoolValue(common.IsUnset(client.Credential)) {
		return common.String(""), nil
	}

	return client.Credential.GetType(), nil
}
func (client *Client) GetBearerToken() (_result *string, _err error) {
	if common.BoolValue(common.IsUnset(client.Credential)) {
		return common.String(""), _err
	}

	_result = client.Credential.GetBearerToken()
	return
}

/**
 * Get accesskey id by using credential
 * @return accesskey id
 */
func (client *Client) GetAccessKeyId() (_result *string, _err error) {
	if common.BoolValue(common.IsUnset(client.Credential)) {
		_result = common.String("")
		return _result, _err
	}

	accessKeyId, _err := client.Credential.GetAccessKeyId()
	if _err != nil {
		return _result, _err
	}

	_result = accessKeyId
	return _result, _err
}

/**
 * Get accesskey secret by using credential
 * @return accesskey secret
 */
func (client *Client) GetAccessKeySecret() (*string, error) {
	if common.BoolValue(common.IsUnset(client.Credential)) {
		return common.String(""), nil
	}

	secret, _err := client.Credential.GetAccessKeySecret()
	if _err != nil {
		return common.String(""), _err
	}
	return secret, _err
}

/**
 * Get security token by using credential
 * @return security token
 */
func (client *Client) GetSecurityToken() (*string, error) {
	if common.BoolValue(common.IsUnset(client.Credential)) {
		return common.String(""), nil
	}

	token, _err := client.Credential.GetSecurityToken()
	if _err != nil {
		return nil, _err
	}
	return token, _err
}

func DefaultAny(inputValue interface{}, defaultValue interface{}) (_result interface{}) {
	if common.BoolValue(common.IsUnset(inputValue)) {
		_result = defaultValue
		return _result
	}

	_result = inputValue
	return _result
}
