package credential

import (
	"os"

	"github.com/lufeijun/go-tool-aliyun/common"
)

// Config is important when call NewCredential
type Config struct {
	Type                  *string  `json:"type"`
	AccessKeyId           *string  `json:"access_key_id"`
	AccessKeySecret       *string  `json:"access_key_secret"`
	OIDCProviderArn       *string  `json:"oidc_provider_arn"`
	OIDCTokenFilePath     *string  `json:"oidc_token"`
	RoleArn               *string  `json:"role_arn"`
	RoleSessionName       *string  `json:"role_session_name"`
	PublicKeyId           *string  `json:"public_key_id"`
	RoleName              *string  `json:"role_name"`
	SessionExpiration     *int     `json:"session_expiration"`
	PrivateKeyFile        *string  `json:"private_key_file"`
	BearerToken           *string  `json:"bearer_token"`
	SecurityToken         *string  `json:"security_token"`
	RoleSessionExpiration *int     `json:"role_session_expiratioon"`
	Policy                *string  `json:"policy"`
	Host                  *string  `json:"host"`
	Timeout               *int     `json:"timeout"`
	ConnectTimeout        *int     `json:"connect_timeout"`
	Proxy                 *string  `json:"proxy"`
	InAdvanceScale        *float64 `json:"inAdvanceScale"`
	Url                   *string  `json:"url"`
	STSEndpoint           *string  `json:"sts_endpoint"`
	ExternalId            *string  `json:"external_id"`
}

func (s Config) String() string {
	return common.Prettify(s)
}

func (s Config) GoString() string {
	return s.String()
}

func (s *Config) SetAccessKeyId(v string) *Config {
	s.AccessKeyId = &v
	return s
}

func (s *Config) SetAccessKeySecret(v string) *Config {
	s.AccessKeySecret = &v
	return s
}

func (s *Config) SetSecurityToken(v string) *Config {
	s.SecurityToken = &v
	return s
}

func (s *Config) SetRoleArn(v string) *Config {
	s.RoleArn = &v
	return s
}

func (s *Config) SetRoleSessionName(v string) *Config {
	s.RoleSessionName = &v
	return s
}

func (s *Config) SetPublicKeyId(v string) *Config {
	s.PublicKeyId = &v
	return s
}

func (s *Config) SetRoleName(v string) *Config {
	s.RoleName = &v
	return s
}

func (s *Config) SetSessionExpiration(v int) *Config {
	s.SessionExpiration = &v
	return s
}

func (s *Config) SetPrivateKeyFile(v string) *Config {
	s.PrivateKeyFile = &v
	return s
}

func (s *Config) SetBearerToken(v string) *Config {
	s.BearerToken = &v
	return s
}

func (s *Config) SetRoleSessionExpiration(v int) *Config {
	s.RoleSessionExpiration = &v
	return s
}

func (s *Config) SetPolicy(v string) *Config {
	s.Policy = &v
	return s
}

func (s *Config) SetHost(v string) *Config {
	s.Host = &v
	return s
}

func (s *Config) SetTimeout(v int) *Config {
	s.Timeout = &v
	return s
}

func (s *Config) SetConnectTimeout(v int) *Config {
	s.ConnectTimeout = &v
	return s
}

func (s *Config) SetProxy(v string) *Config {
	s.Proxy = &v
	return s
}

func (s *Config) SetType(v string) *Config {
	s.Type = &v
	return s
}

func (s *Config) SetOIDCTokenFilePath(v string) *Config {
	s.OIDCTokenFilePath = &v
	return s
}

func (s *Config) SetOIDCProviderArn(v string) *Config {
	s.OIDCProviderArn = &v
	return s
}

func (s *Config) SetURLCredential(v string) *Config {
	if v == "" {
		v = os.Getenv("ALIBABA_CLOUD_CREDENTIALS_URI")
	}
	s.Url = &v
	return s
}

func (s *Config) SetSTSEndpoint(v string) *Config {
	s.STSEndpoint = &v
	return s
}
