// Copyright 2018 The Tigo Authors. All rights reserved.
package WebFramework

import (
	"net/http"
	"fmt"
)

// web容器
type Application struct {
	IPAddress  string
	Port       string
	UrlPattern UrlPattern
}

// 服务启动函数
func (application *Application)Run() {
	application.UrlPattern.Init()
	address := fmt.Sprintf("%s:%s", application.IPAddress, application.Port)
	httpServerErr := http.ListenAndServe(address, nil)
	if httpServerErr != nil {
	}
}
