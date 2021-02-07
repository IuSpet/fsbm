package util

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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

func MinInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func MinInt64(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

func GenerateRandCode(n int) string {
	rand.Seed(time.Now().Unix())
	r := rand.Int63()
	hash := Sha256(strconv.FormatInt(r, 10))
	return hash[:MinInt(32, n)]
}
