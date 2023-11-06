package common

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// SDKError struct is used save error code and message
type SDKError struct {
	Code               *string
	StatusCode         *int
	Message            *string
	Data               *string
	Stack              *string
	errMsg             *string
	Description        *string
	AccessDeniedDetail map[string]interface{}
}

// Set ErrMsg by msg
func (err *SDKError) SetErrMsg(msg string) {
	err.errMsg = String(msg)
}

func (err *SDKError) Error() string {
	if err.errMsg == nil {
		str := fmt.Sprintf("SDKError:\n   StatusCode: %d\n   Code: %s\n   Message: %s\n   Data: %s\n",
			IntValue(err.StatusCode), StringValue(err.Code), StringValue(err.Message), StringValue(err.Data))
		err.SetErrMsg(str)
	}
	return StringValue(err.errMsg)
}

func NewSDKError(obj map[string]interface{}) (err *SDKError) {

	// code 字段转换
	if val, ok := obj["code"].(int); ok {
		err.Code = String(strconv.Itoa(val))
	} else if val, ok := obj["code"].(string); ok {
		err.Code = String(val)
	}

	if obj["message"] != nil {
		err.Message = String(obj["message"].(string))
	}
	if obj["description"] != nil {
		err.Description = String(obj["description"].(string))
	}

	// map 处理
	if detail := obj["accessDeniedDetail"]; detail != nil {
		r := reflect.ValueOf(detail)
		if r.Kind().String() == "map" {
			res := make(map[string]interface{})
			temp := r.MapKeys()
			for _, key := range temp {
				res[key.String()] = r.MapIndex(key).Interface()
			}
			err.AccessDeniedDetail = res
		}
	}

	if data := obj["data"]; data != nil {
		r := reflect.ValueOf(data)
		if r.Kind().String() == "map" {
			res := make(map[string]interface{})
			temp := r.MapKeys()
			for _, key := range temp {
				res[key.String()] = r.MapIndex(key).Interface()
			}

			//
			if statusCode, ok := res["statusCode"]; ok {
				if code, ok := statusCode.(int); ok {
					err.StatusCode = Int(code)
				} else if tmp, ok := statusCode.(string); ok {
					code, err1 := strconv.Atoi(tmp)
					if err1 == nil {
						err.StatusCode = Int(code)
					}
				} else if code, ok := statusCode.(*int); ok {
					err.StatusCode = code
				}
			}

		}

		// 两次转换，性能没问题么
		byt, _ := json.Marshal(data)
		err.Data = String(string(byt))
	}

	if statusCode, ok := obj["statusCode"].(int); ok {
		err.StatusCode = Int(statusCode)
	} else if status, ok := obj["statusCode"].(string); ok {
		statusCode, err_ := strconv.Atoi(status)
		if err_ == nil {
			err.StatusCode = Int(statusCode)
		}
	}

	return
}
