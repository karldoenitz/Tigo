// Copyright 2018 The Tigo Authors. All rights reserved.
package WebFramework

import "net/http"

type UrlPattern struct {
	UrlMapping map[string] interface{Handle(http.ResponseWriter, *http.Request)}
}

// 向http服务挂载单个handler，注意：
//   - handler必须有一个Handle(http.ResponseWriter, *http.Request)函数
func (urlPattern *UrlPattern)AppendUrlPattern(uri string, v interface{Handle(http.ResponseWriter, *http.Request)}) {
	http.HandleFunc(uri, v.Handle)
}

// 初始化url映射，遍历UrlMapping，将handler与对应的URL依次挂载到http服务上
func (urlPattern *UrlPattern)Init() {
	for key, value := range urlPattern.UrlMapping {
		urlPattern.AppendUrlPattern(key, value)
	}
}
