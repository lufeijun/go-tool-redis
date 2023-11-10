package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var basicTypes = []string{
	"int", "int16", "int64", "int32", "float32", "float64", "string", "bool", "uint64", "uint32", "uint16",
}

// Verify whether the parameters meet the requirements
var validateParams = []string{"require", "pattern", "maxLength", "minLength", "maximum", "minimum", "maxItems", "minItems"}

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

func ToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

func DefaultString(reaStr, defaultStr *string) *string {
	if reaStr == nil {
		return defaultStr
	}
	return reaStr
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

func DefaultNumber(reaNum, defaultNum *int) *int {
	if reaNum == nil {
		return defaultNum
	}
	return reaNum
}

// interface 相关

func TransInterfaceToInt(val interface{}) *int {
	if val == nil {
		return nil
	}
	return Int(val.(int))
}

func TransInterfaceToBool(val interface{}) *bool {
	if val == nil {
		return nil
	}

	return Bool(val.(bool))
}

func TransInterfaceToString(val interface{}) *string {
	if val == nil {
		return nil
	}

	return String(val.(string))
}

// help 函数
var jsonParser jsoniter.API

func Convert(in interface{}, out interface{}) error {
	byt, _ := json.Marshal(in)
	decoder := jsonParser.NewDecoder(bytes.NewReader(byt))
	decoder.UseNumber()
	err := decoder.Decode(&out)
	return err
}

func ToReader(obj interface{}) io.Reader {
	switch obj.(type) {
	case *string:
		tmp := obj.(*string)
		return strings.NewReader(StringValue(tmp))
	case []byte:
		return strings.NewReader(string(obj.([]byte)))
	case io.Reader:
		return obj.(io.Reader)
	default:
		panic("Invalid Body. Please set a valid Body.")
	}
}

func Merge(args ...interface{}) map[string]*string {
	finalArg := make(map[string]*string)
	for _, obj := range args {
		switch obj.(type) {
		case map[string]*string:
			arg := obj.(map[string]*string)
			for key, value := range arg {
				if value != nil {
					finalArg[key] = value
				}
			}
		default:
			byt, _ := json.Marshal(obj)
			arg := make(map[string]string)
			err := json.Unmarshal(byt, &arg)
			if err != nil {
				return finalArg
			}
			for key, value := range arg {
				if value != "" {
					finalArg[key] = String(value)
				}
			}
		}
	}

	return finalArg
}

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

func EqualString(val1, val2 *string) *bool {
	return Bool(StringValue(val1) == StringValue(val2))
}

func Validate(params interface{}) error {
	if params == nil {
		return nil
	}
	requestValue := reflect.ValueOf(params)
	if requestValue.IsNil() {
		return nil
	}
	err := validate(requestValue.Elem())
	return err
}
func validate(dataValue reflect.Value) error {
	if strings.HasPrefix(dataValue.Type().String(), "*") { // Determines whether the input is a structure object or a pointer object
		if dataValue.IsNil() {
			return nil
		}
		dataValue = dataValue.Elem()
	}
	dataType := dataValue.Type()
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		valueField := dataValue.Field(i)
		for _, value := range validateParams {
			err := validateParam(field, valueField, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func validateParam(field reflect.StructField, valueField reflect.Value, tagName string) error {
	tag, containsTag := field.Tag.Lookup(tagName) // Take out the checked regular expression
	if containsTag && tagName == "require" {
		err := checkRequire(field, valueField)
		if err != nil {
			return err
		}
	}
	if strings.HasPrefix(field.Type.String(), "[]") { // Verify the parameters of the array type
		err := validateSlice(field, valueField, containsTag, tag, tagName)
		if err != nil {
			return err
		}
	} else if valueField.Kind() == reflect.Ptr { // Determines whether it is a pointer object
		err := validatePtr(field, valueField, containsTag, tag, tagName)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkRequire(field reflect.StructField, valueField reflect.Value) error {
	name, _ := field.Tag.Lookup("json")
	strs := strings.Split(name, ",")
	name = strs[0]
	if !valueField.IsNil() && valueField.IsValid() {
		return nil
	}
	return errors.New(name + " 是必填项")
}

func ValidateModel(a interface{}) error {
	if a == nil {
		return nil
	}
	err := Validate(a)
	return err
}

func validateSlice(field reflect.StructField, valueField reflect.Value, containsregexpTag bool, tag, tagName string) error {
	if valueField.IsValid() && !valueField.IsNil() { // Determines whether the parameter has a value
		if containsregexpTag {
			if tagName == "maxItems" {
				err := checkMaxItems(field, valueField, tag)
				if err != nil {
					return err
				}
			}

			if tagName == "minItems" {
				err := checkMinItems(field, valueField, tag)
				if err != nil {
					return err
				}
			}
		}

		for m := 0; m < valueField.Len(); m++ {
			elementValue := valueField.Index(m)
			if elementValue.Type().Kind() == reflect.Ptr { // Determines whether the child elements of an array are of a basic type
				err := validatePtr(field, elementValue, containsregexpTag, tag, tagName)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func validatePtr(field reflect.StructField, elementValue reflect.Value, containsregexpTag bool, tag, tagName string) error {
	if elementValue.IsNil() {
		return nil
	}
	if isFilterType(elementValue.Elem().Type().String(), basicTypes) {
		if containsregexpTag {
			if tagName == "pattern" {
				err := checkPattern(field, elementValue.Elem(), tag)
				if err != nil {
					return err
				}
			}

			if tagName == "maxLength" {
				err := checkMaxLength(field, elementValue.Elem(), tag)
				if err != nil {
					return err
				}
			}

			if tagName == "minLength" {
				err := checkMinLength(field, elementValue.Elem(), tag)
				if err != nil {
					return err
				}
			}

			if tagName == "maximum" {
				err := checkMaximum(field, elementValue.Elem(), tag)
				if err != nil {
					return err
				}
			}

			if tagName == "minimum" {
				err := checkMinimum(field, elementValue.Elem(), tag)
				if err != nil {
					return err
				}
			}
		}
	} else {
		err := validate(elementValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func isFilterType(realType string, filterTypes []string) bool {
	for _, value := range filterTypes {
		if value == realType {
			return true
		}
	}
	return false
}

// 正则匹配
func checkPattern(field reflect.StructField, valueField reflect.Value, tag string) error {
	if valueField.IsValid() && valueField.String() != "" {
		value := valueField.String()
		r, _ := regexp.Compile("^" + tag + "$")
		if match := r.MatchString(value); !match { // Determines whether the parameter value satisfies the regular expression or not, and throws an error
			return errors.New(value + " is not matched " + tag)
		}
	}
	return nil
}

func checkMaxLength(field reflect.StructField, valueField reflect.Value, tag string) error {
	if valueField.IsValid() && valueField.String() != "" {
		maxLength, err := strconv.Atoi(tag)
		if err != nil {
			return err
		}
		length := valueField.Len()
		if valueField.Kind().String() == "string" {
			length = strings.Count(valueField.String(), "") - 1
		}
		if maxLength < length {
			errMsg := fmt.Sprintf("The length of %s is %d which is more than %d", field.Name, length, maxLength)
			return errors.New(errMsg)
		}
	}
	return nil
}

func checkMinLength(field reflect.StructField, valueField reflect.Value, tag string) error {
	if valueField.IsValid() {
		minLength, err := strconv.Atoi(tag)
		if err != nil {
			return err
		}
		length := valueField.Len()
		if valueField.Kind().String() == "string" {
			length = strings.Count(valueField.String(), "") - 1
		}
		if minLength > length {
			errMsg := fmt.Sprintf("The length of %s is %d which is less than %d", field.Name, length, minLength)
			return errors.New(errMsg)
		}
	}
	return nil
}

func checkMaximum(field reflect.StructField, valueField reflect.Value, tag string) error {
	if valueField.IsValid() && valueField.String() != "" {
		maximum, err := strconv.ParseFloat(tag, 64)
		if err != nil {
			return err
		}
		byt, _ := json.Marshal(valueField.Interface())
		num, err := strconv.ParseFloat(string(byt), 64)
		if err != nil {
			return err
		}
		if maximum < num {
			errMsg := fmt.Sprintf("The size of %s is %f which is greater than %f", field.Name, num, maximum)
			return errors.New(errMsg)
		}
	}
	return nil
}

func checkMinimum(field reflect.StructField, valueField reflect.Value, tag string) error {
	if valueField.IsValid() && valueField.String() != "" {
		minimum, err := strconv.ParseFloat(tag, 64)
		if err != nil {
			return err
		}

		byt, _ := json.Marshal(valueField.Interface())
		num, err := strconv.ParseFloat(string(byt), 64)
		if err != nil {
			return err
		}
		if minimum > num {
			errMsg := fmt.Sprintf("The size of %s is %f which is less than %f", field.Name, num, minimum)
			return errors.New(errMsg)
		}
	}
	return nil
}

func checkMaxItems(field reflect.StructField, valueField reflect.Value, tag string) error {
	if valueField.IsValid() && valueField.String() != "" {
		maxItems, err := strconv.Atoi(tag)
		if err != nil {
			return err
		}
		length := valueField.Len()
		if maxItems < length {
			errMsg := fmt.Sprintf("The length of %s is %d which is more than %d", field.Name, length, maxItems)
			return errors.New(errMsg)
		}
	}
	return nil
}
func checkMinItems(field reflect.StructField, valueField reflect.Value, tag string) error {
	if valueField.IsValid() {
		minItems, err := strconv.Atoi(tag)
		if err != nil {
			return err
		}
		length := valueField.Len()
		if minItems > length {
			errMsg := fmt.Sprintf("The length of %s is %d which is less than %d", field.Name, length, minItems)
			return errors.New(errMsg)
		}
	}
	return nil
}

func isNil(a interface{}) bool {
	defer func() {
		recover()
	}()
	vi := reflect.ValueOf(a)
	return vi.IsNil()
}

func ToMap(args ...interface{}) map[string]interface{} {
	isNotNil := false
	finalArg := make(map[string]interface{})
	for _, obj := range args {
		if obj == nil {
			continue
		}

		if isNil(obj) {
			continue
		}
		isNotNil = true

		switch obj.(type) {
		case map[string]*string:
			arg := obj.(map[string]*string)
			for key, value := range arg {
				if value != nil {
					finalArg[key] = StringValue(value)
				}
			}
		case map[string]interface{}:
			arg := obj.(map[string]interface{})
			for key, value := range arg {
				if value != nil {
					finalArg[key] = value
				}
			}
		case *string:
			str := obj.(*string)
			arg := make(map[string]interface{})
			err := json.Unmarshal([]byte(StringValue(str)), &arg)
			if err == nil {
				for key, value := range arg {
					if value != nil {
						finalArg[key] = value
					}
				}
			}
			tmp := make(map[string]string)
			err = json.Unmarshal([]byte(StringValue(str)), &tmp)
			if err == nil {
				for key, value := range arg {
					if value != "" {
						finalArg[key] = value
					}
				}
			}
		case []byte:
			byt := obj.([]byte)
			arg := make(map[string]interface{})
			err := json.Unmarshal(byt, &arg)
			if err == nil {
				for key, value := range arg {
					if value != nil {
						finalArg[key] = value
					}
				}
				break
			}
		default:
			val := reflect.ValueOf(obj)
			res := structToMap(val)
			for key, value := range res {
				if value != nil {
					finalArg[key] = value
				}
			}
		}
	}

	if !isNotNil {
		return nil
	}
	return finalArg
}

func structToMap(dataValue reflect.Value) map[string]interface{} {
	out := make(map[string]interface{})
	if !dataValue.IsValid() {
		return out
	}
	if dataValue.Kind().String() == "ptr" {
		if dataValue.IsNil() {
			return out
		}
		dataValue = dataValue.Elem()
	}
	if !dataValue.IsValid() {
		return out
	}
	dataType := dataValue.Type()
	if dataType.Kind().String() != "struct" {
		return out
	}
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		name, containsNameTag := field.Tag.Lookup("json")
		if !containsNameTag {
			name = field.Name
		} else {
			strs := strings.Split(name, ",")
			name = strs[0]
		}
		fieldValue := dataValue.FieldByName(field.Name)
		if !fieldValue.IsValid() || fieldValue.IsNil() {
			continue
		}
		if field.Type.String() == "io.Reader" || field.Type.String() == "io.Writer" {
			continue
		} else if field.Type.Kind().String() == "struct" {
			out[name] = structToMap(fieldValue)
		} else if field.Type.Kind().String() == "ptr" &&
			field.Type.Elem().Kind().String() == "struct" {
			if fieldValue.Elem().IsValid() {
				out[name] = structToMap(fieldValue)
			}
		} else if field.Type.Kind().String() == "ptr" {
			if fieldValue.IsValid() && !fieldValue.IsNil() {
				out[name] = fieldValue.Elem().Interface()
			}
		} else if field.Type.Kind().String() == "slice" {
			tmp := make([]interface{}, 0)
			num := fieldValue.Len()
			for i := 0; i < num; i++ {
				value := fieldValue.Index(i)
				if !value.IsValid() {
					continue
				}
				if value.Type().Kind().String() == "ptr" &&
					value.Type().Elem().Kind().String() == "struct" {
					if value.IsValid() && !value.IsNil() {
						tmp = append(tmp, structToMap(value))
					}
				} else if value.Type().Kind().String() == "struct" {
					tmp = append(tmp, structToMap(value))
				} else if value.Type().Kind().String() == "ptr" {
					if value.IsValid() && !value.IsNil() {
						tmp = append(tmp, value.Elem().Interface())
					}
				} else {
					tmp = append(tmp, value.Interface())
				}
			}
			if len(tmp) > 0 {
				out[name] = tmp
			}
		} else {
			out[name] = fieldValue.Interface()
		}

	}
	return out
}

// 请求重试逻辑
func AllowRetry(retry interface{}, retryTimes *int) *bool {
	if IntValue(retryTimes) == 0 {
		return Bool(true)
	}
	retryMap, ok := retry.(map[string]interface{})
	if !ok {
		return Bool(false)
	}
	retryable, ok := retryMap["retryable"].(bool)
	if !ok || !retryable {
		return Bool(false)
	}

	maxAttempts, ok := retryMap["maxAttempts"].(int)
	if !ok || maxAttempts < IntValue(retryTimes) {
		return Bool(false)
	}
	return Bool(true)
}

func GetBackoffTime(backoff interface{}, retrytimes *int) *int {
	backoffMap, ok := backoff.(map[string]interface{})
	if !ok {
		return Int(0)
	}
	policy, ok := backoffMap["policy"].(string)
	if !ok || policy == "no" {
		return Int(0)
	}

	period, ok := backoffMap["period"].(int)
	if !ok || period == 0 {
		return Int(0)
	}

	// 1,2,4,8,16
	maxTime := math.Pow(2.0, float64(IntValue(retrytimes)))
	return Int(rand.Intn(int(maxTime-1)) * period)
}

// 根据返回的 code 判断是否需要重试，ex 500 的话，服务器错误，就不要再重试了
func Retryable(err error) *bool {
	if err == nil {
		return Bool(false)
	}
	if realErr, ok := err.(*SDKError); ok {
		if realErr.StatusCode == nil {
			return Bool(false)
		}
		code := IntValue(realErr.StatusCode)
		return Bool(code >= http.StatusInternalServerError)
	}
	return Bool(true)
}

func Sleep(backoffTime *int) {
	sleeptime := time.Duration(IntValue(backoffTime)) * time.Second
	time.Sleep(sleeptime)
}
