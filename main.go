package main

import (
	"limiter-breaker/limiter"
	"limiter-breaker/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func use(c *gin.Context) {
	c.JSON(200, "nihao")
}

func main() {
	r := gin.Default()
	r.Use(middleware.Limiter(limiter.NewLimiter(time.Second, 4)))
	r.GET("/ping", use)
	r.Run(":8080")
}
