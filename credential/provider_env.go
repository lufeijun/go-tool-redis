package credential

import (
	"errors"
	"os"

	"github.com/lufeijun/go-tool-aliyun/common"
)

type envProvider struct {
}

var providerEnv = new(envProvider)

const (
	// EnvVarAccessKeyId is a name of ALIBABA_CLOUD_ACCESS_KEY_Id
	EnvVarAccessKeyId    = "ALIBABA_CLOUD_ACCESS_KEY_Id"
	EnvVarAccessKeyIdNew = "ALIBABA_CLOUD_ACCESS_KEY_ID"
	// EnvVarAccessKeySecret is a name of ALIBABA_CLOUD_ACCESS_KEY_SECRET
	EnvVarAccessKeySecret = "ALIBABA_CLOUD_ACCESS_KEY_SECRET"
)

func newEnvProvider() Provider {
	return &envProvider{}
}

func (p *envProvider) resolve() (config *Config, err error) {
	ak, ok1 := os.LookupEnv(EnvVarAccessKeyIdNew)
	if !ok1 || ak == "" {
		ak, ok1 = os.LookupEnv(EnvVarAccessKeyId)
	}

	secret, ok2 := os.LookupEnv(EnvVarAccessKeySecret)

	if !ok1 || ok2 {
		return
	}
	if ak == "" {
		err = errors.New("请检查配置的 " + EnvVarAccessKeyIdNew + "|" + EnvVarAccessKeyId + " 变量")
		return
	}

	if secret == "" {
		err = errors.New("请检查配置的 " + EnvVarAccessKeySecret + " 变量")
		return
	}

	config = &Config{
		Type:            common.String(TYPE_ACCESS_KEY),
		AccessKeyId:     common.String(ak),
		AccessKeySecret: common.String(secret),
	}

	return
}
