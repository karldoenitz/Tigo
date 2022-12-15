// Package TigoWeb
// TigoWeb 包是Tigo框架的基础包，内部封装了Handler、http context、middleware、application、session等相关模块的操作。
//
// 使用Tigo创建的服务，可以裸起，也可以通过endless、overseer等进行平滑启动。
// 使用endless平滑启动的示例如下
//
// Basic Example:
//
//	func main() {
//		application := TigoWeb.Application{UrlPatterns: urlRouter}
//		application.EndlessStart()
//	}
//
// 使用overseer平滑启动示例如下，`fetcher.File`是你的Tigo项目的二进制可执行文件的路径，Overseer按照Interval设置的时间轮询该文件是否更新，
// 更新后会进行平滑重启。
//
// Basic Example:
//
//	func main() {
//      application := TigoWeb.Application{UrlPatterns: urlRouter}
//      application.OverseerStart(&fetcher.File{
//  		Path:     "path/to/your/app-file",
//  		Interval: 1 * time.Second,
//  	})
//	}
//
package TigoWeb
