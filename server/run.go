package server

import "github.com/gin-gonic/gin"

func Run() {
	router := gin.Default()
	Register(router)
	_ = router.Run()
}
