package TigoWeb

import (
	"errors"
	"github.com/karldoenitz/Tigo/logger"
	"net/http"
	"time"
)

// TODO 该模块将在v1.6.5中进行大改

// Middleware http中间件
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// chainMiddleware 是http中间件生成器
func chainMiddleware(mw ...Middleware) Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

// InternalServerErrorMiddleware 用来处理控制层出现的异常的中间件
func InternalServerErrorMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("unknown error")
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}

// HttpContextLogMiddleware 记录一个http请求响应时间的中间件
func HttpContextLogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		requestMethod := r.Method
		url := r.RequestURI
		httpResponseWriter := HttpResponseWriter{w, 200}
		defer func() {
			status := httpResponseWriter.GetStatus()
			duration := time.Now().Sub(startTime).Seconds() * 1e3
			switch status {
			case http.StatusInternalServerError:
				logger.Error.Printf("%s | %fms | %s %s", logger.StatusColor(status), duration, requestMethod, url)
				break
			default:
				logger.Info.Printf("%s | %fms | %s %s", logger.StatusColor(status), duration, requestMethod, url)
			}
		}()
		next.ServeHTTP(&httpResponseWriter, r)
	}
}
