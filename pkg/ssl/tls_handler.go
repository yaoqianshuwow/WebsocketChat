package ssl

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"WebsocketChat/pkg/zlog"
	"strconv"
)

func TlsHandler(host string, port int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 配置secure中间件，禁用HTTPS重定向
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect:          false,
			SSLForceHost:         false,
			SSLHost:              host + ":" + strconv.Itoa(port),
			STSSeconds:           0,
			STSIncludeSubdomains: false,
			STSPreload:           false,
			FrameDeny:            false,
			ContentTypeNosniff:   true,
			BrowserXssFilter:     true,
			ContentSecurityPolicy: "default-src 'self'",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// 如果有错误，仅记录警告而不是终止程序
		if err != nil {
			zlog.Warn("Secure middleware error: " + err.Error())
			// 继续处理请求，不中断
		}

		c.Next()
	}
}
