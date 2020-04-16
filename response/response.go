package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Message struct {
	Message string `json:"message"`
}

func Success(ctx *gin.Context, message string) {
	msg := Message{
		Message: message,
	}

	ctx.JSON(http.StatusOK, msg)
}

func Fatal(ctx *gin.Context, message string) {
	msg := Message{
		Message: message,
	}

	ctx.JSON(http.StatusInternalServerError, msg)
}
