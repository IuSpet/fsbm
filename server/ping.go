package server

import "github.com/gin-gonic/gin"

func pong(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"pmsg": "pong"})
}
