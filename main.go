package main

import (
	"errors"
	"limiter-breaker/breaker"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func use(c *gin.Context) {
	c.JSON(200, "nihao")
}

func main() {
	r := gin.Default()
	//r.Use(middleware.Limiter(limiter.NewLimiter(time.Second, 4)))
	r.GET("/ping", use)

	b := breaker.NewBreaker(4, 4, 2, time.Second*30)

	r.GET("/ping1", func(c *gin.Context) {
		err := b.Exec(func() error {
			value, _ := c.GetQuery("value")
			if value == "a" {
				return errors.New("value为a,返回错误")
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}
