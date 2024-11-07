package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.StaticFile("/style.css", "./css/style.css")
	r.StaticFile("/icon.png", "./static/icon.png")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"text": time.Now(),
		})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/:candidate", func(c *gin.Context) {
		name := c.Param("candidate")
		c.HTML(http.StatusOK, "candidate.tmpl", gin.H{
			"name": name,
		})
	})
	r.Run()
}
