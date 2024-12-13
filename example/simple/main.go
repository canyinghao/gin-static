package main

import (
	"log"

	static "github.com/canyinghao/gin-static"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// if Allow DirectoryIndex
	// r.Use(static.Serve("/", static.LocalFile("./public", true)))
	// set prefix
	// r.Use(static.Serve("/static", static.LocalFile("./public", true)))

	r.Use(static.Serve("/", static.LocalFile("./public", false)))

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "test")
	})

	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
