// Copyright 2018 The Tigo Authors. All rights reserved.
package TigoWeb

import (
	"net/http"
	"fmt"
	"Tigo/logger"
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
	logger.Info.Printf("Server run on: %s", address)
	httpServerErr := http.ListenAndServe(address, nil)
	if httpServerErr != nil {
		logger.Error.Printf("HTTP SERVER ERROR! MSG: %s", httpServerErr.Error())
	}
}
