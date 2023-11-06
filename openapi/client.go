package openapi

import "github.com/lufeijun/go-tool-aliyun/common"

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
	// Credential           credential.Credential
	SignatureVersion   *string
	SignatureAlgorithm *string
	Headers            map[string]*string
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

		// credentialConfig := &credential.Config{
		// 	AccessKeyId:     config.AccessKeyId,
		// 	Type:            config.Type,
		// 	AccessKeySecret: config.AccessKeySecret,
		// }
		// credentialConfig.SecurityToken = config.SecurityToken
		// client.Credential, _err = credential.NewCredential(credentialConfig)
		// if _err != nil {
		// 	return _err
		// }

	}
	// else if !tea.BoolValue(util.IsUnset(config.Credential)) {
	// 	client.Credential = config.Credential
	// }

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
