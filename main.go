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
	h := struct {
		Hash     string
		Projects []struct {
			Name        string
			Description string
			Language    string
			Link        string
		}
	}{
		Hash: getBuildHash(),
		Projects: []struct {
			Name        string
			Description string
			Language    string
			Link        string
		}{
			{Name: "fun-banking", Description: "An innovative online banking simulator, used by thousands, designed to provide an engaging and informative platform for individuals to learn about banking.", Language: "Go", Link: "https://fun-banking.com"},
			{Name: "storage-bin", Description: "A storage-like interface using IndexedDB under the hood. It is async, stores any type, and supports large datasets. LocalStorage and SessionStorage behaviors.", Language: "JavaScript", Link: "https://www.npmjs.com/package/@marcellodotgg/storage-bin"},
			{Name: "retroboard-org", Description: "A free online retrospective tool to help lead the team to continuous improvement and stay on track. See feedback in real time, comment, and vote on cards.", Language: "TypeScript", Link: "https://retroboard.org"},
			{Name: "go-estimate", Description: "A straight-to-the-point estimation tool for any sized team. No ads, no limits, just get to guessing in seconds.", Language: "Go", Link: "https://estimate.marcello.gg"},
		},
	}

	router.GET("", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "homepage/index", h)
	})

	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "homepage/index", h)
	})
}

func getBuildHash() string {
	if os.Getenv("GIN_MODE") == "release" {
		return os.Getenv("BUILD_HASH")
	}
	return "local"
}
