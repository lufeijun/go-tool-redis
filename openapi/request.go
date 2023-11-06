package openapi

import (
	"io"

	"github.com/lufeijun/go-tool-aliyun/common"
)

type OpenApiRequest struct {
	Headers          map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	Query            map[string]*string `json:"query,omitempty" xml:"query,omitempty"`
	Body             interface{}        `json:"body,omitempty" xml:"body,omitempty"`
	Stream           io.Reader          `json:"stream,omitempty" xml:"stream,omitempty"`
	HostMap          map[string]*string `json:"hostMap,omitempty" xml:"hostMap,omitempty"`
	EndpointOverride *string            `json:"endpointOverride,omitempty" xml:"endpointOverride,omitempty"`
}

func (s OpenApiRequest) String() string {
	return common.Prettify(s)
}

func (s OpenApiRequest) GoString() string {
	return s.String()
}

func (s *OpenApiRequest) SetHeaders(headers map[string]*string) *OpenApiRequest {
	s.Headers = headers
	return s
}

func (s *OpenApiRequest) SetQuery(queries map[string]*string) *OpenApiRequest {
	s.Query = queries
	return s
}

func (s *OpenApiRequest) SetBody(body interface{}) *OpenApiRequest {
	s.Body = body
	return s
}

func (s *OpenApiRequest) SetStream(stream io.Reader) *OpenApiRequest {
	s.Stream = stream
	return s
}

func (s *OpenApiRequest) SetHostMap(hostMap map[string]*string) *OpenApiRequest {
	s.HostMap = hostMap
	return s
}

func (s *OpenApiRequest) SetEndpointOverride(endpointOverride string) *OpenApiRequest {
	s.EndpointOverride = &endpointOverride
	return s
}

type Params struct {
	Action      *string `json:"action,omitempty" xml:"action,omitempty" require:"true"`
	Version     *string `json:"version,omitempty" xml:"version,omitempty" require:"true"`
	Protocol    *string `json:"protocol,omitempty" xml:"protocol,omitempty" require:"true"`
	Pathname    *string `json:"pathname,omitempty" xml:"pathname,omitempty" require:"true"`
	Method      *string `json:"method,omitempty" xml:"method,omitempty" require:"true"`
	AuthType    *string `json:"authType,omitempty" xml:"authType,omitempty" require:"true"`
	BodyType    *string `json:"bodyType,omitempty" xml:"bodyType,omitempty" require:"true"`
	ReqBodyType *string `json:"reqBodyType,omitempty" xml:"reqBodyType,omitempty" require:"true"`
	Style       *string `json:"style,omitempty" xml:"style,omitempty"`
}

func (s Params) String() string {
	return common.Prettify(s)
}

func (s Params) GoString() string {
	return s.String()
}

func (s *Params) SetAction(action string) *Params {
	s.Action = &action
	return s
}

func (s *Params) SetVersion(version string) *Params {
	s.Version = &version
	return s
}

func (s *Params) SetProtocol(protocol string) *Params {
	s.Protocol = &protocol
	return s
}

func (s *Params) SetPathname(pathname string) *Params {
	s.Pathname = &pathname
	return s
}

func (s *Params) SetMethod(method string) *Params {
	s.Method = &method
	return s
}

func (s *Params) SetAuthType(authType string) *Params {
	s.AuthType = &authType
	return s
}

func (s *Params) SetBodyType(bodyType string) *Params {
	s.BodyType = &bodyType
	return s
}

func (s *Params) SetReqBodyType(reqBodyType string) *Params {
	s.ReqBodyType = &reqBodyType
	return s
}

func (s *Params) SetStyle(style string) *Params {
	s.Style = &style
	return s
}
