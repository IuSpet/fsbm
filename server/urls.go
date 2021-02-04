package server

import "github.com/gin-gonic/gin"

func Register(router *gin.Engine) {
	router.GET("/ping", pong)
}
