package credential

import "github.com/lufeijun/go-tool-aliyun/common"

// AccessKeyCredential is a kind of credential
type AccessKeyCredential struct {
	AccessKeyId     string
	AccessKeySecret string
}

func newAccessKeyCredential(accessKeyId, accessKeySecret string) *AccessKeyCredential {
	return &AccessKeyCredential{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
}

func (s *AccessKeyCredential) GetCredential() (*CredentialModel, error) {
	credential := &CredentialModel{
		AccessKeyId:     common.String(s.AccessKeyId),
		AccessKeySecret: common.String(s.AccessKeySecret),
		Type:            common.String(TYPE_ACCESS_KEY),
	}
	return credential, nil
}

// GetAccessKeyId reutrns  AccessKeyCreential's AccessKeyId
func (a *AccessKeyCredential) GetAccessKeyId() (*string, error) {
	return common.String(a.AccessKeyId), nil
}

// GetAccessSecret reutrns  AccessKeyCreential's AccessKeySecret
func (a *AccessKeyCredential) GetAccessKeySecret() (*string, error) {
	return common.String(a.AccessKeySecret), nil
}

// GetSecurityToken is useless for AccessKeyCreential
func (a *AccessKeyCredential) GetSecurityToken() (*string, error) {
	return common.String(""), nil
}

// GetBearerToken is useless for AccessKeyCreential
func (a *AccessKeyCredential) GetBearerToken() *string {
	return common.String("")
}

// GetType reutrns  AccessKeyCreential's type
func (a *AccessKeyCredential) GetType() *string {
	return common.String(TYPE_ACCESS_KEY)
}
