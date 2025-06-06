package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", PingHandler)
	r.POST("/register", RegisterUser)
}
