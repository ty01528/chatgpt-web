package main

import (
	"chatgpt-web/src"
	"embed"
	"github.com/gin-gonic/gin"
)

//go:embed static/*.html static/css/* static/js/*
var FS embed.FS

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	src.SetRoute(r, FS)
}
