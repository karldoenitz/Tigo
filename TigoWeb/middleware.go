package TigoWeb

import "net/http"

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
