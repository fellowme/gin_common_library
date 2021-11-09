package main

import "github.com/gin-gonic/gin"

func main() {
	app := gin.New()
	app.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
