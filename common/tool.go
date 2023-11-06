package common

import (
	"encoding/json"
	"reflect"
)

// string 相关

// 格式化输出 json 数据
func Prettify(i interface{}) string {
	resp, _ := json.MarshalIndent(i, "", "   ")
	return string(resp)
}

func String(a string) *string {
	return &a
}
func StringValue(a *string) string {
	if a == nil {
		return ""
	}
	return *a
}

// bool 相关
func BoolValue(a *bool) bool {
	if a == nil {
		return false
	}
	return *a
}
func Bool(a bool) *bool {
	return &a
}

// int 相关
func IntValue(a *int) int {
	if a == nil {
		return 0
	}
	return *a
}

func Int(a int) *int {
	return &a
}

// help 函数
func IsUnset(val interface{}) *bool {
	if val == nil {
		return Bool(true)
	}

	// 指针类型
	v := reflect.ValueOf(val)
	kind := v.Kind()
	if kind == reflect.Ptr || kind == reflect.Slice || kind == reflect.Map {
		return Bool(v.IsNil())
	}

	// 零值判断
	valType := reflect.TypeOf(val)
	valZero := reflect.Zero(valType)
	return Bool(valZero == v)
}

func Empty(val *string) *bool {
	return Bool(val == nil || StringValue(val) == "")
}
