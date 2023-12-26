package server

import "github.com/gin-gonic/gin"

//Handling pages with gin router
func routerHandle(router *gin.Engine) {
	router.GET("/", homePage())
}

