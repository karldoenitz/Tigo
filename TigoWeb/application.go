// Package TigoWeb 框架的基本功能包，此包包含了搭建服务的基础功能
package TigoWeb

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gorilla/mux"
	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	"github.com/karldoenitz/Tigo/logger"
	"net/http"
)

// Application web容器
type Application struct {
	IPAddress   string      // IP地址
	Port        int         // 端口
	UrlPatterns []Pattern   // url路由配置
	ConfigPath  string      // 全局配置
	muxRouter   *mux.Router // gorilla的路由
}

// http服务启动函数
func (application *Application) run() {
	address := fmt.Sprintf("%s:%d", application.IPAddress, application.Port)
	var httpServerErr error
	switch {
	// 获取证书与密钥，判断是否启动https服务
	case globalConfig != nil && globalConfig.Cert != "" && globalConfig.CertKey != "":
		logger.Info.Printf("Server run on: https://%s", address)
		httpServerErr = http.ListenAndServeTLS(address, globalConfig.Cert, globalConfig.CertKey, application.muxRouter)
	default:
		logger.Info.Printf("Server run on: http://%s", address)
		httpServerErr = http.ListenAndServe(address, application.muxRouter)
	}
	if httpServerErr != nil {
		logger.Error.Printf("HTTP SERVER ERROR! MSG: %s", httpServerErr.Error())
	}
}

// Listen 端口监听
func (application *Application) Listen(port int) {
	application.Port = port
}

// StartSession 设置session，此函数只提供session操作的接口，以便于第三方session插件嵌入。
func (application *Application) StartSession(sessionInterface SessionInterface, sessionCookieName string) {
	GlobalSessionManager = sessionInterface.NewSessionManager()
	if sessionCookieName != "" {
		SessionCookieName = sessionCookieName
	}
}

// MountFileServer 挂载文件服务
//   - dir 本地文件地址
//   - uris 需要挂载的URI，只支持至多一个URI，输入多个则只取第一个，默认为/路径，URI尽量以/结尾，这样兼容性高一些，比如：
//     application.MountFileServer("/path/to/files", "/files/", "/", "/your/uri/")
func (application *Application) MountFileServer(dir string, uris ...string) {
	if len(uris) == 0 {
		uris = append(uris, "/")
	}
	application.UrlPatterns = append(application.UrlPatterns, Pattern{Url: uris[0], Handler: dir})
}

// Run 服务启动函数
func (application *Application) Run() {
	application.muxRouter = mux.NewRouter()
	application.InitApp()
	application.run()
}

// InitApp 初始化配置信息及路由
func (application *Application) InitApp() {
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
	urlPattern := UrlPattern{UrlPatterns: application.UrlPatterns, router: application.muxRouter}
	urlPattern.Init()
}

// EndlessStart 使用endless进行平滑启动
func (application *Application) EndlessStart() {
	application.InitApp()
	logger.Info.Println("start with endless...")
	err := endless.ListenAndServe(fmt.Sprintf("%s:%d", application.IPAddress, application.Port), application.muxRouter)
	if err != nil {
		panic(fmt.Sprintf("server err: %v", err))
	}
}

// OverseerStart 使用overseer进行平滑启动
//   - fc: overseer包中的fetcher接口，包括file、http、GitHub等
func (application *Application) OverseerStart(fc fetcher.Interface) {
	overseer.Run(overseer.Config{
		Program: application.overseerProgram,
		Address: fmt.Sprintf("%s:%d", application.IPAddress, application.Port),
		Fetcher: fc,
	})
}

func (application *Application) overseerProgram(state overseer.State) {
	logger.Info.Printf("app (%s) start with overseer...\n", state.ID)
	application.InitApp()
	if err := http.Serve(state.Listener, application.muxRouter); err != nil {
		panic(fmt.Sprintf("server err: %v", err))
	}
}
