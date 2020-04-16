package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/goft-cloud/http-proxy/log"
	"github.com/goft-cloud/http-proxy/response"
)

func Recovery(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			message := err.(string)
			log.Error(context, message)
			response.Fatal(context, "internal error")
		}
	}()
	context.Next()
}
