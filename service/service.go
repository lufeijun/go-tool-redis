package service

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/lufeijun/go-tool-aliyun/common"
)

var defaultUserAgent = fmt.Sprintf("AlibabaCloud (%s; %s) Golang/%s Core/%s commonDSL/1", runtime.GOOS, runtime.GOARCH, strings.Trim(runtime.Version(), "go"), "0.01")

type RuntimeOptions struct {
	Autoretry      *bool   `json:"autoretry" xml:"autoretry"`
	IgnoreSSL      *bool   `json:"ignoreSSL" xml:"ignoreSSL"`
	Key            *string `json:"key,omitempty" xml:"key,omitempty"`
	Cert           *string `json:"cert,omitempty" xml:"cert,omitempty"`
	Ca             *string `json:"ca,omitempty" xml:"ca,omitempty"`
	MaxAttempts    *int    `json:"maxAttempts" xml:"maxAttempts"`
	BackoffPolicy  *string `json:"backoffPolicy" xml:"backoffPolicy"`
	BackoffPeriod  *int    `json:"backoffPeriod" xml:"backoffPeriod"`
	ReadTimeout    *int    `json:"readTimeout" xml:"readTimeout"`
	ConnectTimeout *int    `json:"connectTimeout" xml:"connectTimeout"`
	LocalAddr      *string `json:"localAddr" xml:"localAddr"`
	HttpProxy      *string `json:"httpProxy" xml:"httpProxy"`
	HttpsProxy     *string `json:"httpsProxy" xml:"httpsProxy"`
	NoProxy        *string `json:"noProxy" xml:"noProxy"`
	MaxIdleConns   *int    `json:"maxIdleConns" xml:"maxIdleConns"`
	Socks5Proxy    *string `json:"socks5Proxy" xml:"socks5Proxy"`
	Socks5NetWork  *string `json:"socks5NetWork" xml:"socks5NetWork"`
	KeepAlive      *bool   `json:"keepAlive" xml:"keepAlive"`
}

var processStartTime int64 = time.Now().UnixNano() / 1e6
var seqId int64 = 0

// 获取 Goroutine ID
func getGID() uint64 {
	// https://blog.sgmansfield.com/2015/12/goroutine-ids/
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
func (s RuntimeOptions) String() string {
	return common.Prettify(s)
}

func (s RuntimeOptions) GoString() string {
	return s.String()
}

func (s *RuntimeOptions) SetAutoretry(v bool) *RuntimeOptions {
	s.Autoretry = &v
	return s
}

func (s *RuntimeOptions) SetIgnoreSSL(v bool) *RuntimeOptions {
	s.IgnoreSSL = &v
	return s
}

func (s *RuntimeOptions) SetKey(v string) *RuntimeOptions {
	s.Key = &v
	return s
}

func (s *RuntimeOptions) SetCert(v string) *RuntimeOptions {
	s.Cert = &v
	return s
}

func (s *RuntimeOptions) SetCa(v string) *RuntimeOptions {
	s.Ca = &v
	return s
}

func (s *RuntimeOptions) SetMaxAttempts(v int) *RuntimeOptions {
	s.MaxAttempts = &v
	return s
}

func (s *RuntimeOptions) SetBackoffPolicy(v string) *RuntimeOptions {
	s.BackoffPolicy = &v
	return s
}

func (s *RuntimeOptions) SetBackoffPeriod(v int) *RuntimeOptions {
	s.BackoffPeriod = &v
	return s
}

func (s *RuntimeOptions) SetReadTimeout(v int) *RuntimeOptions {
	s.ReadTimeout = &v
	return s
}

func (s *RuntimeOptions) SetConnectTimeout(v int) *RuntimeOptions {
	s.ConnectTimeout = &v
	return s
}

func (s *RuntimeOptions) SetHttpProxy(v string) *RuntimeOptions {
	s.HttpProxy = &v
	return s
}

func (s *RuntimeOptions) SetHttpsProxy(v string) *RuntimeOptions {
	s.HttpsProxy = &v
	return s
}

func (s *RuntimeOptions) SetNoProxy(v string) *RuntimeOptions {
	s.NoProxy = &v
	return s
}

func (s *RuntimeOptions) SetMaxIdleConns(v int) *RuntimeOptions {
	s.MaxIdleConns = &v
	return s
}

func (s *RuntimeOptions) SetLocalAddr(v string) *RuntimeOptions {
	s.LocalAddr = &v
	return s
}

func (s *RuntimeOptions) SetSocks5Proxy(v string) *RuntimeOptions {
	s.Socks5Proxy = &v
	return s
}

func (s *RuntimeOptions) SetSocks5NetWork(v string) *RuntimeOptions {
	s.Socks5NetWork = &v
	return s
}

func (s *RuntimeOptions) SetKeepAlive(v bool) *RuntimeOptions {
	s.KeepAlive = &v
	return s
}

func ReadAsString(body io.Reader) (*string, error) {
	byt, err := io.ReadAll(body)

	if err != nil {
		return common.String(""), err
	}

	r, ok := body.(io.ReadCloser)
	if ok {
		r.Close()
	}
	return common.String(string(byt)), nil

}

func StringifyMapValue(a map[string]interface{}) map[string]*string {
	res := make(map[string]*string)

	for key, v := range a {
		res[key] = ToJSONString(v)
	}

	return res
}

func AnyifyMapValue(a map[string]*string) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range a {
		res[k] = common.StringValue(v)
	}

	return res
}

func ToJSONString(a interface{}) *string {
	switch v := a.(type) {
	case *string:
		return v
	case string:
		return common.String(v)
	case []byte:
		return common.String(string(v))
	case io.Reader:
		byt, err := io.ReadAll(v)
		if err != nil {
			return nil
		}
		return common.String(string(byt))
	}

	byt, err := json.Marshal(a)
	if err != nil {
		return nil
	}
	return common.String(string(byt))
}

func ReadAsBytes(body io.Reader) ([]byte, error) {
	byt, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	r, ok := body.(io.ReadCloser)
	if ok {
		r.Close()
	}

	return byt, nil

}

func DefaultString(reaStr, defaultStr *string) *string {
	if reaStr == nil {
		return defaultStr
	}
	return reaStr
}
func DefaultNumber(reaNum, defaultNum *int) *int {
	if reaNum == nil {
		return defaultNum
	}
	return reaNum
}

func ReadAsJSON(body io.Reader) (result interface{}, err error) {
	byt, err := io.ReadAll(body)
	if err != nil {
		return
	}
	if string(byt) == "" {
		return
	}
	r, ok := body.(io.ReadCloser)
	if ok {
		r.Close()
	}
	d := json.NewDecoder(bytes.NewReader(byt))
	d.UseNumber()
	err = d.Decode(&result)

	return
}

func GetNonce() *string {
	routineId := getGID()
	currentTime := time.Now().UnixNano() / 1e6
	seq := atomic.AddInt64(&seqId, 1)
	randNum := rand.Int63()

	msg := fmt.Sprintf("%d-%d-%d-%d-%d", processStartTime, routineId, currentTime, seq, randNum)

	h := md5.New()
	h.Write([]byte(msg))
	ret := hex.EncodeToString(h.Sum(nil))
	return &ret
}

func Empty(val *string) *bool {
	return common.Bool(val == nil || common.StringValue(val) == "")
}

func ValidateModel(a interface{}) error {
	if a == nil {
		return nil
	}
	return common.Validate(a)
}

func EqualString(val1, val2 *string) *bool {
	return common.Bool(common.StringValue(val1) == common.StringValue(val2))
}

func EqualNumber(val1, val2 *int) *bool {
	return common.Bool(common.IntValue(val1) == common.IntValue(val2))
}

func IsUnset(val interface{}) *bool {
	if val == nil {
		return common.Bool(true)
	}
	v := reflect.ValueOf(val)
	k := v.Kind()
	if k == reflect.Ptr || k == reflect.Map || k == reflect.Slice {
		return common.Bool(v.IsNil())
	}

	valtype := reflect.TypeOf(val)
	valzero := reflect.Zero(valtype)
	return common.Bool(v == valzero)

}

func ToBytes(a *string) []byte {
	return []byte(common.StringValue(a))
}

func AssertAsMap(a interface{}) (_result map[string]interface{}, _err error) {
	r := reflect.ValueOf(a)

	if r.Kind().String() != "map" {
		return nil, errors.New("assertAsMap failed, the type is not map")
	}

	res := make(map[string]interface{})
	tmp := r.MapKeys()
	for _, v := range tmp {
		res[v.String()] = r.MapIndex(v).Interface()
	}

	return res, nil

}

func AssertAsNumber(a interface{}) (_result *int, _err error) {
	res := 0
	switch a.(type) {
	case int:
		tmp := a.(int)
		res = tmp
	case *int:
		tmp := a.(*int)
		res = common.IntValue(tmp)
	default:
		return nil, errors.New(fmt.Sprintf("%v is not a int", a))
	}

	return common.Int(res), nil
}

func AssertAsInteger(value interface{}) (_result *int, _err error) {
	res := 0
	switch value.(type) {
	case int:
		tmp := value.(int)
		res = tmp
	case *int:
		tmp := value.(*int)
		res = common.IntValue(tmp)
	default:
		return nil, errors.New(fmt.Sprintf("%v is not a int", value))
	}

	return common.Int(res), nil
}

func AssertAsBoolean(a interface{}) (_result *bool, _err error) {
	res := false
	switch a.(type) {
	case bool:
		tmp := a.(bool)
		res = tmp
	case *bool:
		tmp := a.(*bool)
		res = common.BoolValue(tmp)
	default:
		return nil, errors.New(fmt.Sprintf("%v is not a bool", a))
	}

	return common.Bool(res), nil
}

func AssertAsString(a interface{}) (_result *string, _err error) {
	res := ""
	switch a.(type) {
	case string:
		tmp := a.(string)
		res = tmp
	case *string:
		tmp := a.(*string)
		res = common.StringValue(tmp)
	default:
		return nil, errors.New(fmt.Sprintf("%v is not a string", a))
	}

	return common.String(res), nil
}

func AssertAsBytes(a interface{}) (_result []byte, _err error) {
	res, ok := a.([]byte)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%v is not a []byte", a))
	}
	return res, nil
}

func AssertAsReadable(a interface{}) (_result io.Reader, _err error) {
	res, ok := a.(io.Reader)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%v is not a reader", a))
	}
	return res, nil
}

func AssertAsArray(a interface{}) (_result []interface{}, _err error) {
	r := reflect.ValueOf(a)
	if r.Kind().String() != "array" && r.Kind().String() != "slice" {
		return nil, errors.New(fmt.Sprintf("%v is not a []interface{}", a))
	}
	aLen := r.Len()
	res := make([]interface{}, r.Len())
	for i := 0; i < aLen; i++ {
		res[i] = r.Index(i).Interface()
	}
	return res, nil
}

func ParseJSON(a *string) interface{} {
	mapTmp := make(map[string]interface{})

	d := json.NewDecoder(bytes.NewReader([]byte(common.StringValue(a))))
	d.UseNumber()
	err := d.Decode(&mapTmp)
	if err == nil {
		return mapTmp
	}

	sliceTmp := make([]interface{}, 0)
	d = json.NewDecoder(bytes.NewReader([]byte(common.StringValue(a))))
	d.UseNumber()
	err = d.Decode(&sliceTmp)
	if err == nil {
		return sliceTmp
	}

	if num, err := strconv.Atoi(common.StringValue(a)); err == nil {
		return num
	}

	if num, err := strconv.ParseBool(common.StringValue(a)); err == nil {
		return num
	}

	if num, err := strconv.ParseFloat(common.StringValue(a), 64); err == nil {
		return num
	}

	return nil

}

func ToString(a []byte) *string {
	return common.String(string(a))
}

func ToMap(in interface{}) map[string]interface{} {
	if in == nil {
		return nil
	}
	res := common.ToMap(in)
	return res
}

func ToFormString(a map[string]interface{}) *string {
	if a == nil {
		return common.String("")
	}
	res := ""
	urlEncoder := url.Values{}
	for key, value := range a {
		v := fmt.Sprintf("%v", value)
		urlEncoder.Add(key, v)
	}
	res = urlEncoder.Encode()
	return common.String(res)
}

func GetDateUTCString() *string {
	return common.String(time.Now().UTC().Format(http.TimeFormat))
}

func GetUserAgent(userAgent *string) *string {
	if userAgent != nil && common.StringValue(userAgent) != "" {
		return common.String(defaultUserAgent + " " + common.StringValue(userAgent))
	}
	return common.String(defaultUserAgent)
}

func Is2xx(code *int) *bool {
	tmp := common.IntValue(code)
	return common.Bool(tmp >= 200 && tmp < 300)
}

func Is3xx(code *int) *bool {
	tmp := common.IntValue(code)
	return common.Bool(tmp >= 300 && tmp < 400)
}

func Is4xx(code *int) *bool {
	tmp := common.IntValue(code)
	return common.Bool(tmp >= 400 && tmp < 500)
}

func Is5xx(code *int) *bool {
	tmp := common.IntValue(code)
	return common.Bool(tmp >= 500 && tmp < 600)
}

func Sleep(millisecond *int) error {
	ms := common.IntValue(millisecond)
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return nil
}

func ToArray(in interface{}) []map[string]interface{} {
	if common.BoolValue(IsUnset(in)) {
		return nil
	}

	tmp := make([]map[string]interface{}, 0)
	byt, _ := json.Marshal(in)
	d := json.NewDecoder(bytes.NewReader(byt))
	d.UseNumber()
	err := d.Decode(&tmp)
	if err != nil {
		return nil
	}
	return tmp
}
