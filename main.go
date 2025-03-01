package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var router = gin.Default()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	setupMiddleware()
	setupRoutes()
	router.Run()
}

func setupMiddleware() {
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(cacheStaticFiles)
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static", "public/")
	router.StaticFS("/.well-known/acme-challenge", http.Dir("/var/www/html/.well-known/acme-challenge"))
}

func cacheStaticFiles(ctx *gin.Context) {
	if strings.Contains(ctx.Request.URL.Path, "/static") {
		ctx.Writer.Header().Set("Cache-Control", "public, max-age=31536000")
	}
	ctx.Next()
}

func setupRoutes() {
	h := struct {
		Hash             string
		GitHubProfileURL string
		Projects         []struct {
			Name        string
			Description string
			Language    string
			Link        string
		}
	}{
		Hash:             getBuildHash(),
		GitHubProfileURL: "https://github.com/marcellodotgg",
		Projects: []struct {
			Name        string
			Description string
			Language    string
			Link        string
		}{
			{Name: "fun-banking", Description: "An innovative online banking simulator, used by thousands, designed to provide an engaging and informative platform for individuals to learn about banking.", Language: "Go", Link: "https://fun-banking.com"},
			{Name: "fwdr", Description: "The (self proclaimed) world's most simple and efficient e-mail forwarder for everyone", Language: "Rust", Link: "https://fwdr.dev"},
			{Name: "go-estimate", Description: "A straight-to-the-point estimation tool for any sized team. No ads, no limits, just get to guessing in seconds.", Language: "Go", Link: "https://estimate.marcello.gg"},
			{Name: "ngx-validators", Description: "A library that provides additional Angular Validators, including the ones Angular provides. Useful for validators in one spot. Supporting a wide-variety of situations.", Language: "TypeScript", Link: "https://github.com/marcellodotgg/ngx-validators"},
			{Name: "axum-template", Description: "An axum template repository with sqlx, docker, and OAuth set up", Language: "Rust", Link: "https://github.com/marcellodotgg/axum-template"},
			{Name: "rusty-deck", Description: "My spin on standard-deck (a Go libary) written in Rust instead of Go.", Language: "Rust", Link: "https://github.com/marcellodotgg/rusty-deck"},
			{Name: "storage-bin", Description: "A storage-like interface using IndexedDB under the hood. It is async, stores any type, and supports large datasets. LocalStorage and SessionStorage behaviors.", Language: "JavaScript", Link: "https://www.npmjs.com/package/@marcellodotgg/storage-bin"},
			{Name: "standard-deck", Description: "A library to create, shuffle, and play with a Standard Deck of Playing Cards. Allows you to focus on the game's rules while this library manages the deck itself.", Language: "Go", Link: "https://github.com/marcellodotgg/standard-deck"},
			{Name: "retroboard-org", Description: "A free online retrospective tool to help lead the team to continuous improvement and stay on track. See feedback in real time, comment, and vote on cards.", Language: "TypeScript", Link: "https://retroboard.org"},
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
