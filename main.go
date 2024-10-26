package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func main() {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	setupMiddleware()
	setupRoutes()
	router.Run()
}

func setupMiddleware() {
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static", "public/")
	router.StaticFS("/.well-known/acme-challenge", http.Dir("/var/www/html/.well-known/acme-challenge"))
}

func setupRoutes() {
	h := struct{}{}

	router.GET("", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "homepage/index", h)
	})

	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "homepage/index", h)
	})
}
