package openapi

import (
	"github.com/lufeijun/go-tool-aliyun/common"
	"github.com/lufeijun/go-tool-aliyun/credential"
)

type GlobalParameters struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	Queries map[string]*string `json:"queries,omitempty" xml:"queries,omitempty"`
}

func (s GlobalParameters) String() string {
	return common.Prettify(s)
}

func (s GlobalParameters) GoString() string {
	return s.String()
}

func (s *GlobalParameters) SetHeaders(h map[string]*string) *GlobalParameters {
	s.Headers = h
	return s
}

func (s *GlobalParameters) SetQueries(queries map[string]*string) *GlobalParameters {
	s.Queries = queries
	return s
}

// Model for initing client ，客户端的一些参数
type Config struct {
	// accesskey id
	AccessKeyId *string `json:"accessKeyId,omitempty" xml:"accessKeyId,omitempty"`
	// accesskey secret
	AccessKeySecret *string `json:"accessKeySecret,omitempty" xml:"accessKeySecret,omitempty"`
	// security token
	SecurityToken *string `json:"securityToken,omitempty" xml:"securityToken,omitempty"`
	// http protocol
	Protocol *string `json:"protocol,omitempty" xml:"protocol,omitempty"`
	// http method
	Method *string `json:"method,omitempty" xml:"method,omitempty"`
	// region id
	RegionId *string `json:"regionId,omitempty" xml:"regionId,omitempty"`
	// read timeout
	ReadTimeout *int `json:"readTimeout,omitempty" xml:"readTimeout,omitempty"`
	// connect timeout
	ConnectTimeout *int `json:"connectTimeout,omitempty" xml:"connectTimeout,omitempty"`
	// http proxy
	HttpProxy *string `json:"httpProxy,omitempty" xml:"httpProxy,omitempty"`
	// https proxy
	HttpsProxy *string `json:"httpsProxy,omitempty" xml:"httpsProxy,omitempty"`
	// credential
	Credential credential.Credential `json:"credential,omitempty" xml:"credential,omitempty"`
	// endpoint
	Endpoint *string `json:"endpoint,omitempty" xml:"endpoint,omitempty"`
	// proxy white list
	NoProxy *string `json:"noProxy,omitempty" xml:"noProxy,omitempty"`
	// max idle conns
	MaxIdleConns *int `json:"maxIdleConns,omitempty" xml:"maxIdleConns,omitempty"`
	// network for endpoint
	Network *string `json:"network,omitempty" xml:"network,omitempty"`
	// user agent
	UserAgent *string `json:"userAgent,omitempty" xml:"userAgent,omitempty"`
	// suffix for endpoint
	Suffix *string `json:"suffix,omitempty" xml:"suffix,omitempty"`
	// socks5 proxy
	Socks5Proxy *string `json:"socks5Proxy,omitempty" xml:"socks5Proxy,omitempty"`
	// socks5 network
	Socks5NetWork *string `json:"socks5NetWork,omitempty" xml:"socks5NetWork,omitempty"`
	// endpoint type
	EndpointType *string `json:"endpointType,omitempty" xml:"endpointType,omitempty"`
	// OpenPlatform endpoint
	OpenPlatformEndpoint *string `json:"openPlatformEndpoint,omitempty" xml:"openPlatformEndpoint,omitempty"`
	// Deprecated
	// credential type
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// Signature Version
	SignatureVersion *string `json:"signatureVersion,omitempty" xml:"signatureVersion,omitempty"`
	// Signature Algorithm
	SignatureAlgorithm *string `json:"signatureAlgorithm,omitempty" xml:"signatureAlgorithm,omitempty"`
	// Global Parameters
	GlobalParameters *GlobalParameters `json:"globalParameters,omitempty" xml:"globalParameters,omitempty"`
	// privite key for client certificate
	Key *string `json:"key,omitempty" xml:"key,omitempty"`
	// client certificate
	Cert *string `json:"cert,omitempty" xml:"cert,omitempty"`
	// server certificate
	Ca *string `json:"ca,omitempty" xml:"ca,omitempty"`
}

func (s Config) String() string {
	return common.Prettify(s)
}

func (s Config) GoString() string {
	return s.String()
}

func (s *Config) SetAccessKeyId(accesskey string) *Config {
	s.AccessKeyId = &accesskey
	return s
}

func (s *Config) SetAccessKeySecret(accessKeySecret string) *Config {
	s.AccessKeySecret = &accessKeySecret
	return s
}

func (s *Config) SetSecurityToken(securityToken string) *Config {
	s.SecurityToken = &securityToken
	return s
}

func (s *Config) SetProtocol(protocol string) *Config {
	s.Protocol = &protocol
	return s
}

func (s *Config) SetMethod(method string) *Config {
	s.Method = &method
	return s
}

func (s *Config) SetRegionId(regionId string) *Config {
	s.RegionId = &regionId
	return s
}

func (s *Config) SetReadTimeout(readTimeout int) *Config {
	s.ReadTimeout = &readTimeout
	return s
}

func (s *Config) SetConnectTimeout(connectTimeout int) *Config {
	s.ConnectTimeout = &connectTimeout
	return s
}

func (s *Config) SetHttpProxy(httpProxy string) *Config {
	s.HttpProxy = &httpProxy
	return s
}

func (s *Config) SetHttpsProxy(httpsProxy string) *Config {
	s.HttpsProxy = &httpsProxy
	return s
}

// func (s *Config) SetCredential(v credential.Credential) *Config {
// 	s.Credential = v
// 	return s
// }

func (s *Config) SetEndpoint(endpoint string) *Config {
	s.Endpoint = &endpoint
	return s
}

func (s *Config) SetNoProxy(noProxy string) *Config {
	s.NoProxy = &noProxy
	return s
}

func (s *Config) SetMaxIdleConns(maxIdleConns int) *Config {
	s.MaxIdleConns = &maxIdleConns
	return s
}

func (s *Config) SetNetwork(network string) *Config {
	s.Network = &network
	return s
}

func (s *Config) SetUserAgent(userAgent string) *Config {
	s.UserAgent = &userAgent
	return s
}

func (s *Config) SetSuffix(suffix string) *Config {
	s.Suffix = &suffix
	return s
}

func (s *Config) SetSocks5Proxy(socks5Proxy string) *Config {
	s.Socks5Proxy = &socks5Proxy
	return s
}

func (s *Config) SetSocks5NetWork(socks5NetWork string) *Config {
	s.Socks5NetWork = &socks5NetWork
	return s
}

func (s *Config) SetEndpointType(endpointType string) *Config {
	s.EndpointType = &endpointType
	return s
}

func (s *Config) SetOpenPlatformEndpoint(openPlatformEndpoint string) *Config {
	s.OpenPlatformEndpoint = &openPlatformEndpoint
	return s
}

func (s *Config) SetType(typestr string) *Config {
	s.Type = &typestr
	return s
}

func (s *Config) SetSignatureVersion(signatureVersion string) *Config {
	s.SignatureVersion = &signatureVersion
	return s
}

func (s *Config) SetSignatureAlgorithm(signatureAlgorithm string) *Config {
	s.SignatureAlgorithm = &signatureAlgorithm
	return s
}

func (s *Config) SetGlobalParameters(globalParameters *GlobalParameters) *Config {
	s.GlobalParameters = globalParameters
	return s
}

func (s *Config) SetKey(key string) *Config {
	s.Key = &key
	return s
}

func (s *Config) SetCert(cert string) *Config {
	s.Cert = &cert
	return s
}

func (s *Config) SetCa(ca string) *Config {
	s.Ca = &ca
	return s
}
