package common

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

type DOFunc func(req *http.Request) (*http.Response, error)

var hookDo = func(fn DOFunc) DOFunc {
	return fn
}

// var hookDo = func(fn func(req *http.Request) (*http.Response, error)) func(req *http.Request) (*http.Response, error) {
// 	return fn
// }

type Request struct {
	Protocol *string
	Port     *int
	Method   *string
	Pathname *string
	Domain   *string
	Headers  map[string]*string
	Query    map[string]*string
	Body     io.Reader
}

// Response is use d wrap http response
type Response struct {
	Body          io.ReadCloser
	StatusCode    *int
	StatusMessage *string
	Headers       map[string]*string
}

// RuntimeObject is used for converting http configuration
type RuntimeObject struct {
	IgnoreSSL      *bool   `json:"ignoreSSL" xml:"ignoreSSL"`
	ReadTimeout    *int    `json:"readTimeout" xml:"readTimeout"`
	ConnectTimeout *int    `json:"connectTimeout" xml:"connectTimeout"`
	LocalAddr      *string `json:"localAddr" xml:"localAddr"`
	HttpProxy      *string `json:"httpProxy" xml:"httpProxy"`
	HttpsProxy     *string `json:"httpsProxy" xml:"httpsProxy"`
	NoProxy        *string `json:"noProxy" xml:"noProxy"`
	MaxIdleConns   *int    `json:"maxIdleConns" xml:"maxIdleConns"`
	Key            *string `json:"key" xml:"key"`
	Cert           *string `json:"cert" xml:"cert"`
	CA             *string `json:"ca" xml:"ca"`
	Socks5Proxy    *string `json:"socks5Proxy" xml:"socks5Proxy"`
	Socks5NetWork  *string `json:"socks5NetWork" xml:"socks5NetWork"`
	// Listener       utils.ProgressListener `json:"listener" xml:"listener"`
	// Tracker        *utils.ReaderTracker   `json:"tracker" xml:"tracker"`
	// Logger         *utils.Logger          `json:"logger" xml:"logger"`
}

func DoRequest(request *Request, requestRuntime map[string]interface{}) (response *Response, err error) {
	runtimeObject := NewRuntimeObject(requestRuntime)

	fieldMap := make(map[string]string)
	// 日志部分
	// utils.InitLogMsg(fieldMap)
	// defer func() {
	// 	if runtimeObject.Logger != nil {
	// 		runtimeObject.Logger.PrintLog(fieldMap, err)
	// 	}
	// }()

	// 设置默认值
	if request.Method == nil {
		request.Method = String("GET")
	}

	if request.Protocol == nil {
		request.Protocol = String("http")
	} else {
		request.Protocol = String(strings.ToLower(StringValue(request.Protocol)))
	}

	// 拼接出 url
	requestURL := ""
	request.Domain = request.Headers["host"]
	if request.Port != nil {
		request.Domain = String(fmt.Sprintf("%s:%d", StringValue(request.Domain), IntValue(request.Port)))
	}
	requestURL = fmt.Sprintf("%s://%s%s", StringValue(request.Protocol), StringValue(request.Domain), StringValue(request.Pathname))
	queryParams := request.Query

	q := url.Values{}
	for key, value := range queryParams {
		q.Add(key, StringValue(value))
	}

	querystring := q.Encode()
	if len(querystring) > 0 {
		if strings.Contains(requestURL, "?") {
			requestURL = fmt.Sprintf("%s&%s", requestURL, querystring)
		} else {
			requestURL = fmt.Sprintf("%s?%s", requestURL, querystring)
		}
	}

	// 发送请求
	httpRequest, err := http.NewRequest(StringValue(request.Method), requestURL, request.Body)
	if err != nil {
		return
	}
	httpRequest.Host = StringValue(request.Domain)

	client := getTeaClient(runtimeObject.getClientTag(StringValue(request.Domain)))
	client.Lock()
	if !client.ifInit {
		trans, err := getHttpTransport(request, runtimeObject)
		if err != nil {
			return nil, err
		}
		client.httpClient.Timeout = time.Duration(IntValue(runtimeObject.ReadTimeout)) * time.Millisecond
		client.httpClient.Transport = trans
		client.ifInit = true
	}
	client.Unlock()

	// 获取了 client 了

	// 请求头处理
	for key, value := range request.Headers {
		if value == nil || key == "content-length" {
			continue
		} else if key == "host" {
			httpRequest.Header["Host"] = []string{StringValue(value)}
			delete(httpRequest.Header, "host")
		} else if key == "user-agent" {
			httpRequest.Header["User-Agent"] = []string{StringValue(value)}
			delete(httpRequest.Header, "user-agent")
		} else {
			httpRequest.Header[key] = []string{StringValue(value)}
		}
	}

	// contentlength, _ := strconv.Atoi(StringValue(request.Headers["content-length"]))
	// event := utils.NewProgressEvent(utils.TransferStartedEvent, 0, int64(contentlength), 0)
	// utils.PublishProgress(runtimeObject.Listener, event)
	putMsgToMap(fieldMap, httpRequest)

	startTime := time.Now()
	fieldMap["{start_time}"] = startTime.Format("2006-01-02 15:04:05")

	// 发送请求了
	res, err := hookDo(client.httpClient.Do)(httpRequest)

	fieldMap["{cost}"] = time.Since(startTime).String() // 请求耗时
	// completedBytes := int64(0)
	// if runtimeObject.Tracker != nil {
	// 	completedBytes = runtimeObject.Tracker.CompletedBytes
	// }
	if err != nil {
		// event = utils.NewProgressEvent(utils.TransferFailedEvent, completedBytes, int64(contentlength), 0)
		// utils.PublishProgress(runtimeObject.Listener, event)
		return
	}

	// event = utils.NewProgressEvent(utils.TransferCompletedEvent, completedBytes, int64(contentlength), 0)
	// utils.PublishProgress(runtimeObject.Listener, event)

	response = NewResponse(res)

	fieldMap["{code}"] = strconv.Itoa(res.StatusCode)
	fieldMap["{res_headers}"] = transToString(res.Header)

	for key, value := range res.Header {
		if len(value) != 0 {
			response.Headers[strings.ToLower(key)] = String(value[0])
		}
	}

	return
}

// 并发控制
type teaClient struct {
	sync.Mutex
	httpClient *http.Client
	ifInit     bool
}

var clientPool = &sync.Map{}

func getTeaClient(tag string) *teaClient {
	client, ok := clientPool.Load(tag)
	if client == nil && !ok {
		client = &teaClient{
			httpClient: &http.Client{},
			ifInit:     false,
		}
		clientPool.Store(tag, client)
	}
	return client.(*teaClient)
}

func (r *RuntimeObject) getClientTag(domain string) string {
	return strconv.FormatBool(BoolValue(r.IgnoreSSL)) + strconv.Itoa(IntValue(r.ReadTimeout)) +
		strconv.Itoa(IntValue(r.ConnectTimeout)) + StringValue(r.LocalAddr) + StringValue(r.HttpProxy) +
		StringValue(r.HttpsProxy) + StringValue(r.NoProxy) + StringValue(r.Socks5Proxy) + StringValue(r.Socks5NetWork) + domain
}

// NewRuntimeObject is used for shortly create runtime object
func NewRuntimeObject(runtime map[string]interface{}) *RuntimeObject {
	if runtime == nil {
		return &RuntimeObject{}
	}

	runtimeObject := &RuntimeObject{
		IgnoreSSL:      TransInterfaceToBool(runtime["ignoreSSL"]),
		ReadTimeout:    TransInterfaceToInt(runtime["readTimeout"]),
		ConnectTimeout: TransInterfaceToInt(runtime["connectTimeout"]),
		LocalAddr:      TransInterfaceToString(runtime["localAddr"]),
		HttpProxy:      TransInterfaceToString(runtime["httpProxy"]),
		HttpsProxy:     TransInterfaceToString(runtime["httpsProxy"]),
		NoProxy:        TransInterfaceToString(runtime["noProxy"]),
		MaxIdleConns:   TransInterfaceToInt(runtime["maxIdleConns"]),
		Socks5Proxy:    TransInterfaceToString(runtime["socks5Proxy"]),
		Socks5NetWork:  TransInterfaceToString(runtime["socks5NetWork"]),
		Key:            TransInterfaceToString(runtime["key"]),
		Cert:           TransInterfaceToString(runtime["cert"]),
		CA:             TransInterfaceToString(runtime["ca"]),
	}
	// if runtime["listener"] != nil {
	// 	runtimeObject.Listener = runtime["listener"].(utils.ProgressListener)
	// }
	// if runtime["tracker"] != nil {
	// 	runtimeObject.Tracker = runtime["tracker"].(*utils.ReaderTracker)
	// }
	// if runtime["logger"] != nil {
	// 	runtimeObject.Logger = runtime["logger"].(*utils.Logger)
	// }
	return runtimeObject
}

// Transport类型实现了RoundTripper接口，支持http、https和http/https代理。Transport类型可以缓存连接以在未来重用。
func getHttpTransport(req *Request, runtime *RuntimeObject) (*http.Transport, error) {
	trans := new(http.Transport)
	httpProxy, err := getHttpProxy(StringValue(req.Protocol), StringValue(req.Domain), runtime)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(*req.Protocol) == "https" {
		if BoolValue(runtime.IgnoreSSL) != true {
			trans.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: false,
			}
			// 证书加载，key，cert
			if runtime.Key != nil && runtime.Cert != nil && StringValue(runtime.Key) != "" && StringValue(runtime.Cert) != "" {
				cert, err := tls.X509KeyPair([]byte(StringValue(runtime.Cert)), []byte(StringValue(runtime.Key)))
				if err != nil {
					return nil, err
				}
				trans.TLSClientConfig.Certificates = []tls.Certificate{cert}
			}

			// CA证书加载。ca 叫签发机构
			if runtime.CA != nil && StringValue(runtime.CA) != "" {
				clientCertPool := x509.NewCertPool()
				ok := clientCertPool.AppendCertsFromPEM([]byte(StringValue(runtime.CA)))
				if !ok {
					return nil, errors.New("Failed to parse root certificate")
				}
				trans.TLSClientConfig.RootCAs = clientCertPool
			}
		} else {
			// 忽略证书校验
			trans.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
	}

	// 代理设置
	if httpProxy != nil {
		trans.Proxy = http.ProxyURL(httpProxy)
		if httpProxy.User != nil {
			password, _ := httpProxy.User.Password()
			auth := httpProxy.User.Username() + ":" + password
			basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
			req.Headers["Proxy-Authorization"] = String(basic)
		}
	}

	//
	if runtime.Socks5Proxy != nil && StringValue(runtime.Socks5Proxy) != "" {
		socks5Proxy, err := getSocks5Proxy(runtime)
		if err != nil {
			return nil, err
		}
		if socks5Proxy != nil {
			var auth *proxy.Auth
			if socks5Proxy.User != nil {
				password, _ := socks5Proxy.User.Password()
				auth = &proxy.Auth{
					User:     socks5Proxy.User.Username(),
					Password: password,
				}
			}
			dialer, err := proxy.SOCKS5(strings.ToLower(StringValue(runtime.Socks5NetWork)), socks5Proxy.String(), auth,
				&net.Dialer{
					Timeout:   time.Duration(IntValue(runtime.ConnectTimeout)) * time.Millisecond,
					DualStack: true,
					LocalAddr: getLocalAddr(StringValue(runtime.LocalAddr)),
				})
			if err != nil {
				return nil, err
			}
			trans.Dial = dialer.Dial
		}
	} else {
		trans.DialContext = setDialContext(runtime)
	}
	if runtime.MaxIdleConns != nil && *runtime.MaxIdleConns > 0 {
		trans.MaxIdleConns = IntValue(runtime.MaxIdleConns)
		trans.MaxIdleConnsPerHost = IntValue(runtime.MaxIdleConns)
	}

	return trans, nil
}

func getHttpProxy(protocol, host string, runtime *RuntimeObject) (proxy *url.URL, err error) {

	// 直接请求
	urls := getNoProxy(protocol, runtime)
	for _, url := range urls {
		if url == host {
			return nil, nil
		}
	}

	// 拿代理连接
	if protocol == "https" {
		if runtime.HttpsProxy != nil && StringValue(runtime.HttpsProxy) != "" {
			proxy, err = url.Parse(StringValue(runtime.HttpsProxy))
		} else if rawurl := os.Getenv("HTTPS_PROXY"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		} else if rawurl := os.Getenv("https_proxy"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		}
	} else {
		if runtime.HttpProxy != nil && StringValue(runtime.HttpProxy) != "" {
			proxy, err = url.Parse(StringValue(runtime.HttpProxy))
		} else if rawurl := os.Getenv("HTTP_PROXY"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		} else if rawurl := os.Getenv("http_proxy"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		}
	}

	return proxy, err
}

func getNoProxy(protocol string, runtime *RuntimeObject) []string {
	var urls []string

	if runtime.NoProxy != nil && StringValue(runtime.NoProxy) != "" {
		urls = strings.Split(StringValue(runtime.NoProxy), ",")
	} else if rawurl := os.Getenv("NO_PROXY"); rawurl != "" {
		urls = strings.Split(rawurl, ",")
	} else if rawurl := os.Getenv("no_proxy"); rawurl != "" {
		urls = strings.Split(rawurl, ",")
	}
	return urls
}

func getSocks5Proxy(runtime *RuntimeObject) (proxy *url.URL, err error) {
	if runtime.Socks5Proxy != nil && StringValue(runtime.Socks5Proxy) != "" {
		proxy, err = url.Parse(StringValue(runtime.Socks5Proxy))
	}
	return proxy, err
}

func getLocalAddr(localAddr string) (addr *net.TCPAddr) {
	if localAddr != "" {
		addr = &net.TCPAddr{
			IP: []byte(localAddr),
		}
	}
	return addr
}

func setDialContext(runtime *RuntimeObject) func(cxt context.Context, net, addr string) (c net.Conn, err error) {
	return func(cxt context.Context, network, addr string) (c net.Conn, err error) {
		if runtime.LocalAddr != nil && StringValue(runtime.LocalAddr) != "" {
			netAddr := &net.TCPAddr{
				IP: []byte(StringValue(runtime.LocalAddr)),
			}
			return (&net.Dialer{
				Timeout:   time.Duration(IntValue(runtime.ConnectTimeout)) * time.Second,
				DualStack: true,
				LocalAddr: netAddr,
			}).DialContext(cxt, network, addr)
		}

		return (&net.Dialer{
			Timeout:   time.Duration(IntValue(runtime.ConnectTimeout)) * time.Second,
			DualStack: true,
		}).DialContext(cxt, network, addr)

	}
}

// NewRequest is used shortly create Request
func NewRequest() (req *Request) {
	return &Request{
		Headers: map[string]*string{},
		Query:   map[string]*string{},
	}
}

// NewResponse is create response with http response
func NewResponse(httpResponse *http.Response) (res *Response) {
	res = &Response{}
	res.Body = httpResponse.Body
	res.Headers = make(map[string]*string)
	res.StatusCode = Int(httpResponse.StatusCode)
	res.StatusMessage = String(httpResponse.Status)
	return
}

func putMsgToMap(fieldMap map[string]string, request *http.Request) {
	fieldMap["{host}"] = request.Host
	fieldMap["{method}"] = request.Method
	fieldMap["{uri}"] = request.URL.RequestURI()
	fieldMap["{pid}"] = strconv.Itoa(os.Getpid())
	fieldMap["{version}"] = strings.Split(request.Proto, "/")[1]
	hostname, _ := os.Hostname()
	fieldMap["{hostname}"] = hostname
	fieldMap["{req_headers}"] = transToString(request.Header)
	fieldMap["{target}"] = request.URL.Path + request.URL.RawQuery
}

func transToString(object interface{}) string {
	byt, _ := json.Marshal(object)
	return string(byt)
}
