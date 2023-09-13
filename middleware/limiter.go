package middleware

import (
	"limiter-breaker/limiter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Limiter(l *limiter.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !l.Allow() {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "服务繁忙",
			})
			ctx.Abort()
		}
		ctx.Next()
	}

}
