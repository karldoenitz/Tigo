// Package web
// web 包是Tigo框架的基础包，内部封装了Handler、http context、middleware、application、session等相关模块的操作。
//
// 使用Tigo创建的服务，可以裸起，也可以通过endless、overseer等进行平滑启动。
// 使用endless平滑启动的示例如下
//
// Basic Example:
//
//	func main() {
//	    application := web.Application{UrlPatterns: urlRouter}
//	    application.EndlessStart()
//	}
//
// 使用overseer平滑启动示例如下，`fetcher.File`是你的Tigo项目的二进制可执行文件的路径，Overseer按照Interval设置的时间轮询该文件是否更新，
// 更新后会进行平滑重启。
//
// Basic Example:
//
//		func main() {
//	     application := web.Application{UrlPatterns: urlRouter}
//	     application.OverseerStart(&fetcher.File{
//	         Path:     "path/to/your/app-file",
//	         Interval: 1 * time.Second,
//	     })
//		}
//
// ---------------------------------------------------------------------------------------------------------------------
//
// 通过中间件可以设置context上下文，在中间件中设置`context.Context`后，可以在handler中获取，只要能获取`http.Request`，就可以从中获取在
// 中间件中设置的`context.Context`。示例如下：
//
// Basic Example:
//
//	func Authorize(w *http.ResponseWriter, r *http.Request) bool {
//		cxt := context.WithValue(r.Context(), "keyFat", "valueJu")
//		*r = *r.WithContext(cxt)
//		return true
//	}
//
//	func (p *PingHandler) Post() {
//		valueInCtx := p.Request.Context().Value("keyFat")
//		p.ResponseAsText(valueInCtx)
//	}
package web
