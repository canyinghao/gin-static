# Gin Static Middleware

[![Run Tests](https://github.com/soulteary/gin-static/actions/workflows/go.yml/badge.svg)](https://github.com/soulteary/gin-static/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/soulteary/static/branch/main/graph/badge.svg)](https://codecov.io/gh/soulteary/static)
[![Go Report Card](https://goreportcard.com/badge/github.com/soulteary/gin-static)](https://goreportcard.com/report/github.com/soulteary/gin-static)
[![GoDoc](https://godoc.org/github.com/soulteary/gin-static?status.svg)](https://godoc.org/github.com/soulteary/gin-static)

Static middleware, support both local files and embed filesystem.

No historical burden, using go 1.21, 100% code coverage.

## Quick Start

### Download and Import

Download and install it:

```bash
go get github.com/soulteary/gin-static
```

Import it in your code:

```go
import "github.com/soulteary/gin-static"
```

## Example

See the [example](example)

### Serve Local Files

[local files]:# (example/simple/main.go)

```go
package main

import (
	"log"

	static "github.com/soulteary/gin-static"
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
```

### Serve Embed folder

[embedmd]:# (example/embed/main.go)

```go
package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	static "github.com/soulteary/gin-static"
	"github.com/gin-gonic/gin"
)

//go:embed public
var EmbedFS embed.FS

func main() {
	r := gin.Default()

	staticFiles, err := static.EmbedFolder(EmbedFS, "public/page")
	if err != nil {
		log.Fatalln("initialization of embed folder failed:", err)
	} else {
		r.Use(static.Serve("/", staticFiles))
	}

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "test")
	})

	r.NoRoute(func(c *gin.Context) {
		fmt.Printf("%s doesn't exists, redirect on /\n", c.Request.URL.Path)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
```
