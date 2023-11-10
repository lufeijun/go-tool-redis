package cdn

import (
	"github.com/lufeijun/go-tool-aliyun/common"
	"github.com/lufeijun/go-tool-aliyun/openapi"
	"github.com/lufeijun/go-tool-aliyun/service"
)

type Client struct {
	openapi.Client
}

func NewClient(config *openapi.Config) (*Client, error) {
	client := new(Client)
	err := client.Init(config)
	return client, err
}

func (client *Client) Init(config *openapi.Config) (_err error) {
	_err = client.Client.Init(config)
	if _err != nil {
		return _err
	}
	client.EndpointRule = common.String("central")
	client.EndpointMap = map[string]*string{
		"ap-northeast-1": common.String("cdn.ap-southeast-1.aliyuncs.com"),
		"ap-south-1":     common.String("cdn.ap-southeast-1.aliyuncs.com"),
		"ap-southeast-1": common.String("cdn.ap-southeast-1.aliyuncs.com"),
		"ap-southeast-2": common.String("cdn.ap-southeast-1.aliyuncs.com"),
		"ap-southeast-3": common.String("cdn.ap-southeast-1.aliyuncs.com"),
		"ap-southeast-5": common.String("cdn.ap-southeast-1.aliyuncs.com"),
		"eu-central-1":   common.String("cdn.ap-southeast-1.aliyuncs.com"),
		"eu-west-1":      common.String("cdn.ap-southeast-1.aliyuncs.com"),
		"me-east-1":      common.String("cdn.ap-southeast-1.aliyuncs.com"),
		"us-east-1":      common.String("cdn.ap-southeast-1.aliyuncs.com"),
		"us-west-1":      common.String("cdn.ap-southeast-1.aliyuncs.com"),
	}
	_err = client.CheckConfig(config)
	if _err != nil {
		return _err
	}
	client.Endpoint, _err = client.GetEndpoint(common.String("cdn"), client.RegionId, client.EndpointRule, client.Network, client.Suffix, client.EndpointMap, client.Endpoint)
	if _err != nil {
		return _err
	}

	return nil
}

func (client *Client) GetEndpoint(productId *string, regionId *string, endpointRule *string, network *string, suffix *string, endpointMap map[string]*string, endpoint *string) (_result *string, _err error) {
	if !common.BoolValue(common.Empty(endpoint)) {
		_result = endpoint
		return _result, _err
	}

	if !common.BoolValue(common.IsUnset(endpointMap)) && !common.BoolValue(common.Empty(endpointMap[common.StringValue(regionId)])) {
		_result = endpointMap[common.StringValue(regionId)]
		return _result, _err
	}

	_body, _err := GetEndpointRules(productId, regionId, endpointRule, network, suffix)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

// 刷新Object缓存

type RefreshObjectCachesRequest struct {
	Force      *bool   `json:"Force,omitempty" xml:"Force,omitempty"`
	ObjectPath *string `json:"ObjectPath,omitempty" xml:"ObjectPath,omitempty"`
	// The type of the object that you want to refresh. Valid values:
	//
	// *   **File** (default): refreshes one or more files.
	// *   **Directory**: refreshes the files in one or more directories.
	// *   **Regex**: refreshes content based on regular expressions.
	// *   **ExQuery**: omits parameters after the question mark in the URL and refreshes content.
	//
	// If you set the ObjectType parameter to File or Directory, you can view [Refresh and prefetch resources](~~27140~~) to obtain more information. If you set the ObjectType parameter to Regex, you can view [Configure URL refresh rules that contain regular expressions](~~146195~~) to obtain more information.
	//
	// If you set the ObjectType parameter to Directory, the resources in the directory that you want to refresh are marked as expired. You cannot delete the directory. If clients request resources on POPs that are marked as expired, Alibaba Cloud CDN checks whether the resources on your origin server are updated. If resources are updated, Alibaba Cloud CDN retrieves the latest version of the resources and returns the resources to the clients. Otherwise, the origin server returns the 304 status code.
	ObjectType    *string `json:"ObjectType,omitempty" xml:"ObjectType,omitempty"`
	OwnerId       *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	SecurityToken *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
}

func (s RefreshObjectCachesRequest) String() string {
	return common.Prettify(s)
}

func (s RefreshObjectCachesRequest) GoString() string {
	return s.String()
}

func (s *RefreshObjectCachesRequest) SetForce(v bool) *RefreshObjectCachesRequest {
	s.Force = &v
	return s
}

func (s *RefreshObjectCachesRequest) SetObjectPath(v string) *RefreshObjectCachesRequest {
	s.ObjectPath = &v
	return s
}

func (s *RefreshObjectCachesRequest) SetObjectType(v string) *RefreshObjectCachesRequest {
	s.ObjectType = &v
	return s
}

func (s *RefreshObjectCachesRequest) SetOwnerId(v int64) *RefreshObjectCachesRequest {
	s.OwnerId = &v
	return s
}

func (s *RefreshObjectCachesRequest) SetSecurityToken(v string) *RefreshObjectCachesRequest {
	s.SecurityToken = &v
	return s
}

type RefreshObjectCachesResponseBody struct {
	// The refresh task ID. If multiple tasks are returned, the IDs are separated by commas (,). The task IDs are merged based on the following rules:
	//
	// *   If the tasks are specified for the same accelerated domain name, submitted within the same second, and run to refresh content based on URLs inscommond of directories, the task IDs are merged into one task ID.
	// *   If the number of tasks that are specified for the same accelerated domain name, submitted within the same second, and run to refresh content based on URLs inscommond of directories exceeds 2,000, every 2,000 task IDs are merged into one task ID.
	RefreshTaskId *string `json:"RefreshTaskId,omitempty" xml:"RefreshTaskId,omitempty"`
	// The request ID.
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s RefreshObjectCachesResponseBody) String() string {
	return common.Prettify(s)
}

func (s RefreshObjectCachesResponseBody) GoString() string {
	return s.String()
}

func (s *RefreshObjectCachesResponseBody) SetRefreshTaskId(v string) *RefreshObjectCachesResponseBody {
	s.RefreshTaskId = &v
	return s
}

func (s *RefreshObjectCachesResponseBody) SetRequestId(v string) *RefreshObjectCachesResponseBody {
	s.RequestId = &v
	return s
}

type RefreshObjectCachesResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *RefreshObjectCachesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s RefreshObjectCachesResponse) String() string {
	return common.Prettify(s)
}

func (s RefreshObjectCachesResponse) GoString() string {
	return s.String()
}

func (s *RefreshObjectCachesResponse) SetHeaders(v map[string]*string) *RefreshObjectCachesResponse {
	s.Headers = v
	return s
}

func (s *RefreshObjectCachesResponse) SetStatusCode(v int32) *RefreshObjectCachesResponse {
	s.StatusCode = &v
	return s
}

func (s *RefreshObjectCachesResponse) SetBody(v *RefreshObjectCachesResponseBody) *RefreshObjectCachesResponse {
	s.Body = v
	return s
}

func (client *Client) RefreshObjectCaches(request *RefreshObjectCachesRequest) (_result *RefreshObjectCachesResponse, _err error) {
	runtime := &service.RuntimeOptions{}
	_result = &RefreshObjectCachesResponse{}
	_body, _err := client.RefreshObjectCachesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) RefreshObjectCachesWithOptions(request *RefreshObjectCachesRequest, runtime *service.RuntimeOptions) (_result *RefreshObjectCachesResponse, _err error) {
	_err = common.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !common.BoolValue(common.IsUnset(request.Force)) {
		query["Force"] = request.Force
	}

	if !common.BoolValue(common.IsUnset(request.ObjectPath)) {
		query["ObjectPath"] = request.ObjectPath
	}

	if !common.BoolValue(common.IsUnset(request.ObjectType)) {
		query["ObjectType"] = request.ObjectType
	}

	if !common.BoolValue(common.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !common.BoolValue(common.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	req := &openapi.OpenApiRequest{
		Query: service.Query(query),
	}
	params := &openapi.Params{
		Action:      common.String("RefreshObjectCaches"),
		Version:     common.String("2018-05-10"),
		Protocol:    common.String("HTTPS"),
		Pathname:    common.String("/"),
		Method:      common.String("POST"),
		AuthType:    common.String("AK"),
		Style:       common.String("RPC"),
		ReqBodyType: common.String("formData"),
		BodyType:    common.String("json"),
	}
	_result = &RefreshObjectCachesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = common.Convert(_body, &_result)
	return _result, _err
}
