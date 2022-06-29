// Package TigoWeb
// TigoWeb 包是Tigo框架的基础包，内部封装了Handler、http context、middleware、application、session等相关模块的操作。
//
// 使用Tigo创建的服务，可以裸起，也可以通过graceful、endless、overseer等进行平滑启动。
// 使用endless平滑启动的示例如下
//
// Basic Example:
//
//	func main() {
//		application := TigoWeb.Application{UrlRouters: urlRouter}
//		application.PrepareStart()
//		endless.DefaultReadTimeOut = 10 * time.Second
//		endless.DefaultWriteTimeOut = 10 * time.Second
//		err := endless.ListenAndServe(fmt.Sprintf(":%d", 4000), http.DefaultServeMux)
//		if err != nil {
//			panic(fmt.Sprintf("server err:%v", err))
//		}
//	}
//
// 使用overseer平滑启动示例如下
//
// Basic Example:
//
//	func main() {
//		overseer.Run(overseer.Config{
//			Program: prog,
//			Address: ":3000",
//			Fetcher: &fetcher.File{
//				Path:     "./test",
//				Interval: 1 * time.Second,
//			},
//		})
//	}
//
//	func prog(state overseer.State) {
//		log.Printf("app (%s) listening...", state.ID)
//		application := TigoWeb.Application{UrlRouters: urlRouter}
//		application.InitApp()
//		http.Serve(state.Listener, nil)
//	}
package TigoWeb
