package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.GetHeader("X-API-KEY") != "DEMO" {
			context.AbortWithStatus(http.StatusUnauthorized)
		}
		context.Next()
	}
}
