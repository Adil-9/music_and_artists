package server

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func StartConfigureServer() {
	router := gin.Default()

	routerHandle(router)

	// Get the absolute path to the templates directory
	templatesDir := filepath.Join(".", "templates")
	staticDir := filepath.Join(".", "static")

	// Serve static files from the specified directory
	router.Static("/static", staticDir)
	// Load templates from the specified directory
	router.LoadHTMLGlob(filepath.Join(templatesDir, "*.html"))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
