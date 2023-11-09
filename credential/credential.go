package credential

import (
	"errors"

	"github.com/lufeijun/go-tool-aliyun/common"
	"github.com/lufeijun/go-tool-aliyun/debug"
)

const (
	TYPE_ACCESS_KEY = "access_key"
)

var debuglog = debug.Init("credential")

var hookParse = func(err error) error {
	return err
}

// 接口定义
// Credential is an interface for getting actual credential
type Credential interface {
	GetAccessKeyId() (*string, error)
	GetAccessKeySecret() (*string, error)
	GetSecurityToken() (*string, error)
	GetBearerToken() *string
	GetType() *string
	GetCredential() (*CredentialModel, error)
}

// cerdentail 的实现
// Environmental virables that may be used by the provider
const (
	ENVCredentialFile  = "ALIBABA_CLOUD_CREDENTIALS_FILE"
	ENVEcsMetadata     = "ALIBABA_CLOUD_ECS_METADATA"
	PATHCredentialFile = "~/.alibabacloud/credentials"
	ENVRoleArn         = "ALIBABA_CLOUD_ROLE_ARN"
	ENVOIDCProviderArn = "ALIBABA_CLOUD_OIDC_PROVIDER_ARN"
	ENVOIDCTokenFile   = "ALIBABA_CLOUD_OIDC_TOKEN_FILE"
	ENVRoleSessionName = "ALIBABA_CLOUD_ROLE_SESSION_NAME"
)

// 默认读取配置的地方
type Provider interface {
	resolve() (*Config, error)
}

type providerChain struct {
	Providers []Provider
}

// var defaultproviders = []Provider{providerEnv, providerOIDC, providerProfile, providerInstance}
var defaultproviders = []Provider{providerEnv}
var defaultChain = newProviderChain(defaultproviders)

func newProviderChain(providers []Provider) Provider {
	return &providerChain{
		Providers: providers,
	}
}

func (p *providerChain) resolve() (*Config, error) {
	for _, provider := range p.Providers {
		config, err := provider.resolve()
		if err != nil {
			return nil, err
		} else if config == nil {
			continue
		}
		return config, err
	}
	return nil, errors.New("No credential found")

}

// 根据 type 不同类型，返回不同的实现
// NewCredential return a credential according to the type in config.
// if config is nil, the function will use default provider chain to get credential.
// please see README.md for detail.
func NewCredential(config *Config) (credential Credential, err error) {
	if config == nil {
		config, err = defaultChain.resolve()
		if err != nil {
			return
		}
		return NewCredential(config)
	}
	switch common.StringValue(config.Type) {
	// case "credentials_uri":
	// 	credential = newURLCredential(tea.StringValue(config.Url))
	// case "oidc_role_arn":
	// 	err = checkoutAssumeRamoidc(config)
	// 	if err != nil {
	// 		return
	// 	}
	// 	runtime := &utils.Runtime{
	// 		Host:           tea.StringValue(config.Host),
	// 		Proxy:          tea.StringValue(config.Proxy),
	// 		ReadTimeout:    tea.IntValue(config.Timeout),
	// 		ConnectTimeout: tea.IntValue(config.ConnectTimeout),
	// 		STSEndpoint:    tea.StringValue(config.STSEndpoint),
	// 	}
	// 	credential = newOIDCRoleArnCredential(tea.StringValue(config.AccessKeyId), tea.StringValue(config.AccessKeySecret), tea.StringValue(config.RoleArn), tea.StringValue(config.OIDCProviderArn), tea.StringValue(config.OIDCTokenFilePath), tea.StringValue(config.RoleSessionName), tea.StringValue(config.Policy), tea.IntValue(config.RoleSessionExpiration), runtime)
	case "access_key":
		err = checkAccessKey(config)
		if err != nil {
			return
		}
		credential = newAccessKeyCredential(common.StringValue(config.AccessKeyId), common.StringValue(config.AccessKeySecret))
	// case "sts":
	// 	err = checkSTS(config)
	// 	if err != nil {
	// 		return
	// 	}
	// 	credential = newStsTokenCredential(tea.StringValue(config.AccessKeyId), tea.StringValue(config.AccessKeySecret), tea.StringValue(config.SecurityToken))
	// case "ecs_ram_role":
	// 	checkEcsRAMRole(config)
	// 	runtime := &utils.Runtime{
	// 		Host:           tea.StringValue(config.Host),
	// 		Proxy:          tea.StringValue(config.Proxy),
	// 		ReadTimeout:    tea.IntValue(config.Timeout),
	// 		ConnectTimeout: tea.IntValue(config.ConnectTimeout),
	// 	}
	// 	credential = newEcsRAMRoleCredential(tea.StringValue(config.RoleName), tea.Float64Value(config.InAdvanceScale), runtime)
	// case "ram_role_arn":
	// 	err = checkRAMRoleArn(config)
	// 	if err != nil {
	// 		return
	// 	}
	// 	runtime := &utils.Runtime{
	// 		Host:           tea.StringValue(config.Host),
	// 		Proxy:          tea.StringValue(config.Proxy),
	// 		ReadTimeout:    tea.IntValue(config.Timeout),
	// 		ConnectTimeout: tea.IntValue(config.ConnectTimeout),
	// 		STSEndpoint:    tea.StringValue(config.STSEndpoint),
	// 	}
	// 	credential = newRAMRoleArnWithExternalIdCredential(
	// 		tea.StringValue(config.AccessKeyId),
	// 		tea.StringValue(config.AccessKeySecret),
	// 		tea.StringValue(config.RoleArn),
	// 		tea.StringValue(config.RoleSessionName),
	// 		tea.StringValue(config.Policy),
	// 		tea.IntValue(config.RoleSessionExpiration),
	// 		tea.StringValue(config.ExternalId),
	// 		runtime)
	// case "rsa_key_pair":
	// 	err = checkRSAKeyPair(config)
	// 	if err != nil {
	// 		return
	// 	}
	// 	file, err1 := os.Open(tea.StringValue(config.PrivateKeyFile))
	// 	if err1 != nil {
	// 		err = fmt.Errorf("InvalidPath: Can not open PrivateKeyFile, err is %s", err1.Error())
	// 		return
	// 	}
	// 	defer file.Close()
	// 	var privateKey string
	// 	scan := bufio.NewScanner(file)
	// 	for scan.Scan() {
	// 		if strings.HasPrefix(scan.Text(), "----") {
	// 			continue
	// 		}
	// 		privateKey += scan.Text() + "\n"
	// 	}
	// 	runtime := &utils.Runtime{
	// 		Host:           tea.StringValue(config.Host),
	// 		Proxy:          tea.StringValue(config.Proxy),
	// 		ReadTimeout:    tea.IntValue(config.Timeout),
	// 		ConnectTimeout: tea.IntValue(config.ConnectTimeout),
	// 		STSEndpoint:    tea.StringValue(config.STSEndpoint),
	// 	}
	// 	credential = newRsaKeyPairCredential(privateKey, tea.StringValue(config.PublicKeyId), tea.IntValue(config.SessionExpiration), runtime)
	// case "bearer":
	// 	if tea.StringValue(config.BearerToken) == "" {
	// 		err = errors.New("BearerToken cannot be empty")
	// 		return
	// 	}
	// 	credential = newBearerTokenCredential(tea.StringValue(config.BearerToken))
	default:
		err = errors.New("Invalid type option, support: access_key, sts, ecs_ram_role, ram_role_arn, rsa_key_pair")
		return
	}
	return credential, nil
}

func checkAccessKey(config *Config) (err error) {
	if common.StringValue(config.AccessKeyId) == "" {
		err = errors.New("AccessKeyId 不能为空")
	} else if common.StringValue(config.AccessKeySecret) == "" {
		err = errors.New("AccessKeySecret 不能为空")
	}

	return
}
