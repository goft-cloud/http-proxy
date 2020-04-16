package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func Request(context *gin.Context) {
	uniqueId := uuid.NewV4().String()
	traceid := context.Request.Header.Get("traceid")
	if traceid == "" {
		traceid = uniqueId
	}

	spanid := context.Request.Header.Get("spanid")
	if spanid == "" {
		spanid = uniqueId
	}

	parentid := context.Request.Header.Get("parentid")
	if parentid == "" {
		parentid = uniqueId
	}

	context.Set("traceid", traceid)
	context.Set("spanid", spanid)
	context.Set("parentid", parentid)

	context.Writer.Header().Set("X-Request-Id", spanid)
	context.Next()
}
