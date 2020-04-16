package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/goft-cloud/http-proxy/log"
)

func Recovery(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			message := err.(string)
			log.Error(context, message)
			context.JSON(500, gin.H{
				"message": "internal error",
			})
		}
	}()
	context.Next()
}
