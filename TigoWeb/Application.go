// Copyright 2018 The Tigo Authors. All rights reserved.
package TigoWeb

import (
	"net/http"
	"fmt"
	"github.com/karldoenitz/Tigo/logger"
)

// web容器
type Application struct {
	IPAddress  string      // IP地址
	Port       string      // 端口
	UrlPattern UrlPattern  // url路由配置
	ConfigPath string      // 全局配置
}

// http服务启动函数
func (application *Application)run() {
	application.UrlPattern.Init()
	address := fmt.Sprintf("%s:%s", application.IPAddress, application.Port)
	logger.Info.Printf("Server run on: %s", address)
	httpServerErr := http.ListenAndServe(address, nil)
	if httpServerErr != nil {
		logger.Error.Printf("HTTP SERVER ERROR! MSG: %s", httpServerErr.Error())
	}
}

// https服务启动函数
func (application *Application)runTLS(cert string, key string) {
	application.UrlPattern.Init()
	address := fmt.Sprintf("%s:%s", application.IPAddress, application.Port)
	logger.Info.Printf("Server run on: %s", address)
	http.ListenAndServeTLS(address, cert, key, nil)
}

// 服务启动函数
func (application *Application)Run() {
	// 初始化全局变量
	InitGlobalConfig(application.ConfigPath)
	// 获取证书与密钥，判断是否启动https服务
	cert, certKey := globalConfig.Cert, globalConfig.CertKey
	if cert != "" && certKey != "" {
		application.runTLS(cert, certKey)
	} else {
		application.run()
	}
}
