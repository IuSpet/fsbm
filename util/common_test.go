package util

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	fmt.Println(Md5("123456"))
	fmt.Println(Md5("123456789"))
	fmt.Println(Md5("654321"))
}

func TestSha256(t *testing.T) {
	fmt.Println(Sha256("123456"))
	fmt.Println(Sha256("123456789"))
	fmt.Println(Sha256("654321"))
}
