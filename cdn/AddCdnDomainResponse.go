package cdn

import "github.com/lufeijun/go-tool-aliyun/common"

type AddCdnDomainResponseBody struct {
	// The ID of the request.
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s AddCdnDomainResponseBody) String() string {
	return common.Prettify(s)
}

func (s AddCdnDomainResponseBody) GoString() string {
	return s.String()
}

func (s *AddCdnDomainResponseBody) SetRequestId(v string) *AddCdnDomainResponseBody {
	s.RequestId = &v
	return s
}

type AddCdnDomainResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *AddCdnDomainResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s AddCdnDomainResponse) String() string {
	return common.Prettify(s)
}

func (s AddCdnDomainResponse) GoString() string {
	return s.String()
}

func (s *AddCdnDomainResponse) SetHeaders(v map[string]*string) *AddCdnDomainResponse {
	s.Headers = v
	return s
}

func (s *AddCdnDomainResponse) SetStatusCode(v int32) *AddCdnDomainResponse {
	s.StatusCode = &v
	return s
}

func (s *AddCdnDomainResponse) SetBody(v *AddCdnDomainResponseBody) *AddCdnDomainResponse {
	s.Body = v
	return s
}
