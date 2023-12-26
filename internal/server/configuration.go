package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func StartConfigureServer() {
	router := gin.Default()

	routerHandle(router)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
