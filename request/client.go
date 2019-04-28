// Package request 提供Tigo框架自带的http client功能，此包包含发送http请求的方法。
package request

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/karldoenitz/Tigo/logger"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// HttpClient 是自定义HTTPClient
type HttpClient struct {
	*http.Client
}

// 执行http请求
func (client HttpClient) request(method, uri string, headers map[string]string, bodyReader io.Reader) (res *Response, err error) {
	res = &Response{}
	// 创建新的请求
	req, err := http.NewRequest(method, uri, bodyReader)
	if err != nil {
		return nil, err
	}
	// 设置请求头
	if host, ok := headers["Host"]; ok {
		req.Host = host
	}
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()
	if response.StatusCode == http.StatusOK {
		switch response.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ := gzip.NewReader(response.Body)
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)

				if err != nil && err != io.EOF {
					logger.Error.Println("Read response error!")
				}

				if n == 0 {
					break
				}
				res.Content = append(res.Content, buf...)
			}
		default:
			res.Content, err = ioutil.ReadAll(response.Body)
		}
	}
	return
}

// Response 自定义Http的Response
type Response struct {
	*http.Response
	Content []byte
}

// ToContentStr 将Response实例的Content转换为字符串
func (response *Response) ToContentStr() string {
	return string(response.Content)
}

// Request 发送指定的Request请求
func Request(method string, requestUrl string, postParams map[string]interface{}, headers ...map[string]string) (*Response, error) {
	client := &HttpClient{http.DefaultClient}
	requestHeaders := map[string]string{}
	if len(headers) > 0 {
		requestHeaders = headers[0]
	}
	if _, ok := requestHeaders["Content-Type"]; !ok {
		requestHeaders["Content-Type"] = "application/x-www-form-urlencoded"
	}
	contentType := requestHeaders["Content-Type"]
	var response *Response
	var err error
	switch {
	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
		postValues := url.Values{}
		for key, value := range postParams {
			postValues.Set(key, value.(string))
		}
		postStr := postValues.Encode()
		response, err = client.request(method, requestUrl, requestHeaders, strings.NewReader(postStr))
		if err != nil {
			logger.Error.Printf("%s %s ERROR\n", method, requestUrl)
			return nil, err
		}
	case strings.HasPrefix(contentType, "application/json"):
		postData, err := json.Marshal(postParams)
		if err != nil {
			logger.Error.Println("marshal to json failed")
			return nil, err
		}
		response, err = client.request(method, requestUrl, requestHeaders, bytes.NewReader(postData))
		if err != nil {
			logger.Error.Printf("%s %s ERROR\n", method, requestUrl)
			return nil, err
		}
	default:
		break
	}
	return response, nil
}

// Get 向指定url发送get请求
func Get(requestUrl string, headers ...map[string]string) (*Response, error) {
	client := &HttpClient{http.DefaultClient}
	var requestHeaders map[string]string
	if len(headers) > 0 {
		requestHeaders = headers[0]
	}
	response, err := client.request("GET", requestUrl, requestHeaders, nil)
	if err != nil {
		logger.Error.Printf("GET %s ERROR\n", requestUrl)
		return nil, err
	}
	return response, nil
}

// Post 向指定url发送post请求
func Post(requestUrl string, postParams map[string]interface{}, headers ...map[string]string) (*Response, error) {
	return Request("POST", requestUrl, postParams, headers...)
}

// Put 向指定url发送put请求
func Put(requestUrl string, postParams map[string]interface{}, headers ...map[string]string) (*Response, error) {
	return Request("PUT", requestUrl, postParams, headers...)
}

// Patch 向指定url发送patch请求
func Patch(requestUrl string, postParams map[string]interface{}, headers ...map[string]string) (*Response, error) {
	return Request("PATCH", requestUrl, postParams, headers...)
}

// Head 向指定url发送head请求
func Head(requestUrl string, headers ...map[string]string) (*Response, error) {
	client := &HttpClient{http.DefaultClient}
	var requestHeaders map[string]string
	if len(headers) > 0 {
		requestHeaders = headers[0]
	}
	response, err := client.request("HEAD", requestUrl, requestHeaders, nil)
	if err != nil {
		logger.Error.Printf("HEAD %s ERROR\n", requestUrl)
		return nil, err
	}
	return response, nil
}

// Options 向指定url发送options请求
func Options(requestUrl string, headers ...map[string]string) (*Response, error) {
	client := &HttpClient{http.DefaultClient}
	var requestHeaders map[string]string
	if len(headers) > 0 {
		requestHeaders = headers[0]
	}
	response, err := client.request("OPTIONS", requestUrl, requestHeaders, nil)
	if err != nil {
		logger.Error.Printf("OPTIONS %s ERROR\n", requestUrl)
		return nil, err
	}
	return response, nil
}

// Delete 向指定url发送delete请求
func Delete(requestUrl string, headers ...map[string]string) (*Response, error) {
	client := &HttpClient{http.DefaultClient}
	var requestHeaders map[string]string
	if len(headers) > 0 {
		requestHeaders = headers[0]
	}
	response, err := client.request("DELETE", requestUrl, requestHeaders, nil)
	if err != nil {
		logger.Error.Printf("DELETE %s ERROR\n", requestUrl)
		return nil, err
	}
	return response, nil
}
