package cdn

import "github.com/lufeijun/go-tool-aliyun/common"

type AddCdnDomainRequest struct {
	// The workload type of the accelerated domain name. Valid values:
	//
	// *   **web**: images and small files
	// *   **download**: large files
	// *   **video**: on-demand video and audio streaming
	CdnType *string `json:"CdnType,omitempty" xml:"CdnType,omitempty"`
	// The URL that is used to check the accessibility of the origin server.
	CheckUrl *string `json:"CheckUrl,omitempty" xml:"CheckUrl,omitempty"`
	// The domain name that you want to add to Alibaba Cloud CDN.
	//
	// A wildcard domain that starts with a period (.) is supported, such as .example.com.
	DomainName   *string `json:"DomainName,omitempty" xml:"DomainName,omitempty"`
	OwnerAccount *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId      *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	// The ID of the resource group.
	//
	// If you do not set this parameter, the system uses the ID of the default resource group.
	ResourceGroupId *string `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
	// The acceleration region. Default value: domestic. Valid values:
	//
	// *   **domestic**: Chinese mainland
	// *   **overseas**: global (excluding the Chinese mainland)
	// *   **global**: global
	Scope         *string `json:"Scope,omitempty" xml:"Scope,omitempty"`
	SecurityToken *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
	// The information about the addresses of origin servers.
	Sources *string `json:"Sources,omitempty" xml:"Sources,omitempty"`
	// Details about the tags. You can specify up to 20 tags.
	Tag []*AddCdnDomainRequestTag `json:"Tag,omitempty" xml:"Tag,omitempty" type:"Repeated"`
	// The top-level domain.
	TopLevelDomain *string `json:"TopLevelDomain,omitempty" xml:"TopLevelDomain,omitempty"`
}

func (s AddCdnDomainRequest) String() string {
	return common.Prettify(s)
}

func (s AddCdnDomainRequest) GoString() string {
	return s.String()
}

func (s *AddCdnDomainRequest) SetCdnType(v string) *AddCdnDomainRequest {
	s.CdnType = &v
	return s
}

func (s *AddCdnDomainRequest) SetCheckUrl(v string) *AddCdnDomainRequest {
	s.CheckUrl = &v
	return s
}

func (s *AddCdnDomainRequest) SetDomainName(v string) *AddCdnDomainRequest {
	s.DomainName = &v
	return s
}

func (s *AddCdnDomainRequest) SetOwnerAccount(v string) *AddCdnDomainRequest {
	s.OwnerAccount = &v
	return s
}

func (s *AddCdnDomainRequest) SetOwnerId(v int64) *AddCdnDomainRequest {
	s.OwnerId = &v
	return s
}

func (s *AddCdnDomainRequest) SetResourceGroupId(v string) *AddCdnDomainRequest {
	s.ResourceGroupId = &v
	return s
}

func (s *AddCdnDomainRequest) SetScope(v string) *AddCdnDomainRequest {
	s.Scope = &v
	return s
}

func (s *AddCdnDomainRequest) SetSecurityToken(v string) *AddCdnDomainRequest {
	s.SecurityToken = &v
	return s
}

func (s *AddCdnDomainRequest) SetSources(v string) *AddCdnDomainRequest {
	s.Sources = &v
	return s
}

func (s *AddCdnDomainRequest) SetTag(v []*AddCdnDomainRequestTag) *AddCdnDomainRequest {
	s.Tag = v
	return s
}

func (s *AddCdnDomainRequest) SetTopLevelDomain(v string) *AddCdnDomainRequest {
	s.TopLevelDomain = &v
	return s
}

type AddCdnDomainRequestTag struct {
	// The key of the tag. Valid values of N: **1 to 20**.
	Key *string `json:"Key,omitempty" xml:"Key,omitempty"`
	// The value of the tag. Valid values of N: **1 to 20**.
	Value *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s AddCdnDomainRequestTag) String() string {
	return common.Prettify(s)
}

func (s AddCdnDomainRequestTag) GoString() string {
	return s.String()
}

func (s *AddCdnDomainRequestTag) SetKey(v string) *AddCdnDomainRequestTag {
	s.Key = &v
	return s
}

func (s *AddCdnDomainRequestTag) SetValue(v string) *AddCdnDomainRequestTag {
	s.Value = &v
	return s
}
