package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/goft-cloud/http-proxy/log"
	"io/ioutil"
	"time"
)

func Log(context *gin.Context) {
	// Start timer
	start := time.Now()
	path := context.Request.URL.String()

	// Post origin body
	body := ""
	if context.Request.Method != "GET" {
		bs, _ := ioutil.ReadAll(context.Request.Body)
		body = string(bs)

		// Reset
		context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bs))

		// parse
		_ = context.Request.ParseForm()
	}

	// Process request
	context.Next()

	end := time.Now()
	cost := end.Sub(start).Nanoseconds()/1000000

	ip := context.ClientIP()
	method := context.Request.Method
	status := context.Writer.Status()
	message := context.Errors.ByType(gin.ErrorTypePrivate).String()

	bodySize := context.Writer.Size()

	log.GetEntry(context).
		WithField("status", status).
		WithField("uri", path).
		WithField("method", method).
		WithField("cost", cost).
		WithField("keys", context.Keys).
		WithField("params", context.Params).
		WithField("header", context.Request.Header).
		WithField("body", body).
		WithField("ip", ip).
		WithField("size", bodySize).
		WithField("postForm", context.Request.PostForm).
		Notice(message)
}

