package util

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Md5(raw string) string {
	h := md5.New()
	_, _ = h.Write([]byte(raw))
	return fmt.Sprintf("%+X", h.Sum(nil))
}

func Sha256(raw string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(raw))
	return fmt.Sprintf("%+X", h.Sum(nil))
}

func ErrorJson(ctx *gin.Context, status int, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"Status":  status,
		"message": msg,
	})
}

func EndJson(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "OK",
		"data":    data,
	})
}
