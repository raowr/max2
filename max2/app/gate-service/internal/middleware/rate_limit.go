package middleware

import (
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"golang.org/x/time/rate"
)

// RateLimit 返回一个限流中间件，每秒允许 limit 个请求，桶容量为 burst
func RateLimit(limit rate.Limit, burst int) ghttp.HandlerFunc {
	limiter := rate.NewLimiter(limit, burst)
	return func(r *ghttp.Request) {
		if !limiter.Allow() {
			r.Response.WriteStatus(http.StatusTooManyRequests)
			r.Response.WriteJson(g.Map{
				"code":    429,
				"message": "Too Many Requests",
			})
			r.ExitAll()
			return
		}
		r.Middleware.Next()
	}
}
