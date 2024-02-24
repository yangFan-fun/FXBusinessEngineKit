package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func LoadTLS() gin.HandlerFunc {
	return func(context *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:8080",
		})
		err := middleware.Process(context.Writer, context.Request)
		if err != nil {
			return
		}
		context.Next()
	}
}
