package middleware

import "github.com/gin-gonic/gin"

// CORS 跨域访问
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		header := c.Writer.Header()
		header.Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		header.Set("Access-Control-Allow-Credentials", "true")
		header.Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
		)
		header.Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

	}
}
