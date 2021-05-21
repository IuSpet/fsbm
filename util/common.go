package util

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"fsbm/util/redis"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const verificationExpiration = 3 * time.Minute

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
		"status": status,
		"msg":   msg,
	})
}

func EndJson(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":   "OK",
		"data":   data,
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

// 产生n位，最多32位随机字符
func GenerateRandCode(n int) string {
	rand.Seed(time.Now().Unix())
	r := rand.Int63()
	hash := Sha256(strconv.FormatInt(r, 10))
	return hash[:MinInt(32, n)]
}

func SetFileTransportHeader(ctx *gin.Context, fileName string) {
	ctx.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(fileName))
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	ctx.Header("Cache-Control", "must-revalidate")
	ctx.Header("Pragma", "public")
	ctx.Header("Access-Control-Expose-Headers", "*")
	ctx.Header("FileName", fileName)
	ctx.File(fileName)
}

func SetImageTransportHeader(ctx *gin.Context) {
	ctx.Header("Content-Type", "image/gif")
}

func SetVerificationCode(ctx context.Context, email string) (string, error) {
	code := GenerateRandCode(6)
	key := fmt.Sprintf(UserVerificationCodeTemplate, email)
	err := redis.SetWithRetry(ctx, key, code, verificationExpiration)
	if err != nil {
		return code, err
	}
	return code, nil
}

func GetVerificationCode(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf(UserVerificationCodeTemplate, email)
	res, err := redis.GetWithRetry(ctx, key)
	if err != nil {
		return "", err
	}
	return res, nil
}

func LikeCondition(field string) string {
	return "%" + field + "%"
}

func CmpInterface(x, y interface{}) bool {
	switch x.(type) {
	case string:
		return x.(string) < y.(string)
	case int:
		return x.(int) < y.(int)
	case int64:
		return x.(int64) < y.(int64)
	case int32:
		return x.(int32) < y.(int32)
	case int16:
		return x.(int16) < y.(int16)
	case int8:
		return x.(int8) < y.(int8)
	case time.Time:
		return x.(time.Time).Before(y.(time.Time))
	}
	return true
}
