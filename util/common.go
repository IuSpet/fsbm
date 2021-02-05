package util

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
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
