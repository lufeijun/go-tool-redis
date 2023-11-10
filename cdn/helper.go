package cdn

import (
	"fmt"
	"strings"

	"github.com/lufeijun/go-tool-aliyun/common"
)

func GetEndpointRules(product, regionId, endpointType, network, suffix *string) (_result *string, _err error) {
	if common.StringValue(endpointType) == "regional" {
		if common.StringValue(regionId) == "" {
			_err = fmt.Errorf("RegionId is empty, please set a valid RegionId")
			return common.String(""), _err
		}
		_result = common.String(strings.Replace("<product><suffix><network>.<region_id>.aliyuncs.com",
			"<region_id>", common.StringValue(regionId), 1))
	} else {
		_result = common.String("<product><suffix><network>.aliyuncs.com")
	}
	_result = common.String(strings.Replace(common.StringValue(_result),
		"<product>", strings.ToLower(common.StringValue(product)), 1))
	if common.StringValue(network) == "" || common.StringValue(network) == "public" {
		_result = common.String(strings.Replace(common.StringValue(_result), "<network>", "", 1))
	} else {
		_result = common.String(strings.Replace(common.StringValue(_result),
			"<network>", "-"+common.StringValue(network), 1))
	}
	if common.StringValue(suffix) == "" {
		_result = common.String(strings.Replace(common.StringValue(_result), "<suffix>", "", 1))
	} else {
		_result = common.String(strings.Replace(common.StringValue(_result),
			"<suffix>", "-"+common.StringValue(suffix), 1))
	}
	return _result, nil
}
