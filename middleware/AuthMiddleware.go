package middleware

import (
	"ddns/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 获取 session
		session := sessions.Default(c)
		// 检查 session 中是否存在用户信息
		user := session.Get("user")
		if user == nil && user != "ld@omyue.com" {
			msg := utils.MesTpl("fail", "未经授权", "")
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
			c.Abort()
			return
		}
		// 将用户信息传递到下一个处理程序
		c.Set("user", user)
		// 继续处理请求
		c.Next()
	}
}
