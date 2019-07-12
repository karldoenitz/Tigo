// Package TigoWeb 框架的基本功能包，此包包含了搭建服务的基础功能
package TigoWeb

import (
	"fmt"
	"github.com/karldoenitz/Tigo/logger"
	"net/http"
)

// Application web容器
type Application struct {
	IPAddress  string                 // IP地址
	Port       int                    // 端口
	UrlPattern map[string]interface{} // url路由配置
	UrlRouters []Router               // url路由配置
	ConfigPath string                 // 全局配置
}

// http服务启动函数
func (application *Application) run() {
	address := fmt.Sprintf("%s:%d", application.IPAddress, application.Port)
	logger.Info.Printf("Server run on: %s", address)
	httpServerErr := http.ListenAndServe(address, nil)
	if httpServerErr != nil {
		logger.Error.Printf("HTTP SERVER ERROR! MSG: %s", httpServerErr.Error())
	}
}

// https服务启动函数
func (application *Application) runTLS(cert string, key string) {
	address := fmt.Sprintf("%s:%d", application.IPAddress, application.Port)
	logger.Info.Printf("Server run on: %s", address)
	http.ListenAndServeTLS(address, cert, key, nil)
}

// Listen 端口监听
func (application *Application) Listen(port int) {
	application.Port = port
}

// MountFileServer 挂载文件服务
//  - dir 本地文件地址
//  - uris 需要挂载的URI列表
func (application *Application) MountFileServer(dir string, uris ...string) {
	if len(uris) > 0 {
		for _, uri := range uris {
			http.Handle(uri, http.FileServer(http.Dir(dir)))
		}
	} else {
		http.Handle("/", http.FileServer(http.Dir(dir)))
	}
}


// Run 服务启动函数
func (application *Application) Run() {
	// 初始化全局变量
	if application.ConfigPath != "" {
		InitGlobalConfig(application.ConfigPath)
	}
	if globalConfig != nil && globalConfig.IP != "" {
		application.IPAddress = globalConfig.IP
	}
	if globalConfig != nil && globalConfig.Port != 0 {
		application.Port = globalConfig.Port
	}
	// url挂载
	urlPattern := UrlPattern{UrlMapping: application.UrlPattern, UrlRouters: application.UrlRouters}
	urlPattern.Init()
	// 获取证书与密钥，判断是否启动https服务
	cert, certKey := "", ""
	if globalConfig != nil {
		cert, certKey = globalConfig.Cert, globalConfig.CertKey
	}
	if cert != "" && certKey != "" {
		application.runTLS(cert, certKey)
	} else {
		application.run()
	}
}
