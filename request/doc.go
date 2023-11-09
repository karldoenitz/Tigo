// Package request
// request 包是Tigo框架的http客户端包，主要用来向服务端发送http请求，该包主要用来做http接口调用，目前尚不支持通过http接口上传文件，但是相
// 关代码可以二次封装，开发者可以自行拓展接口，实现文件上传。
// ---------------------------------------------------------------------------------------------------------------------
// 发送基础的http请求，代码示例如下：
//
// Basic Example:
//
//	url := "https://www.github.com"
//	response, err := request.Get(url)
//	if err != nil {
//		panic(err.Error())
//	}
//	print(response.ToContentStr())
//
// ---------------------------------------------------------------------------------------------------------------------
// 如果要post请求，发送一个json到服务端，可以参考如下示例：
//
// Basic Example:
//
//	url := "https://your.server.address:port/request-url"
//	headers := map[string]string{
//		"Content-Type": "application/json",
//	}
//	param := map[string]interface{}{
//		"param": "value",
//	}
//	response, err := request.Post(url, param, headers)
//	print(response.ToContentStr())
//
// ---------------------------------------------------------------------------------------------------------------------
// `request.Response`继承了`http.Response`，封装了成员变量Content，主要是http response报文的报文体的内容，方便开发者查看报文以及做一
// 些定制化开发。
package request
