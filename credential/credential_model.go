package credential

import "github.com/lufeijun/go-tool-aliyun/common"

// CredentialModel is a model
type CredentialModel struct {
	// accesskey id
	AccessKeyId *string `json:"accessKeyId,omitempty" xml:"accessKeyId,omitempty"`
	// accesskey secret
	AccessKeySecret *string `json:"accessKeySecret,omitempty" xml:"accessKeySecret,omitempty"`
	// security token
	SecurityToken *string `json:"securityToken,omitempty" xml:"securityToken,omitempty"`
	// bearer token
	BearerToken *string `json:"bearerToken,omitempty" xml:"bearerToken,omitempty"`
	// type
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s CredentialModel) String() string {
	return common.Prettify(s)
}

func (s CredentialModel) GoString() string {
	return s.String()
}

func (s *CredentialModel) SetAccessKeyId(accessKeyId string) *CredentialModel {
	s.AccessKeyId = &accessKeyId
	return s
}

func (s *CredentialModel) SetAccessKeySecret(accessKeySecret string) *CredentialModel {
	s.AccessKeySecret = &accessKeySecret
	return s
}

func (s *CredentialModel) SetSecurityToken(securityToken string) *CredentialModel {
	s.SecurityToken = &securityToken
	return s
}

func (s *CredentialModel) SetBearerToken(bearerToken string) *CredentialModel {
	s.BearerToken = &bearerToken
	return s
}

func (s *CredentialModel) SetType(typestr string) *CredentialModel {
	s.Type = &typestr
	return s
}
