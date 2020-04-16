package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goft-cloud/http-proxy/config"
	"github.com/goft-cloud/http-proxy/middleware"
	"github.com/goft-cloud/http-proxy/response"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	s *server
)

func Init() error {
	s = &server{
		gin:    gin.New(),
		isInit: true,
	}

	gin.SetMode(gin.DebugMode)

	// Set mode
	if config.App.IsRelease() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Register middleware and must define by first
	s.gin.Use(middleware.Request, middleware.Log, middleware.Recovery)

	// Add default route
	s.addDefaultRoute()

	// Server address
	s.Addr = fmt.Sprintf("%s:%d", config.App.Host, config.App.Port)
	return nil
}

// Start server
func Run() error {
	service := &http.Server{
		Addr:    s.Addr,
		Handler: s.gin,
	}

	// Start server by coroutine
	go func() {
		if err := service.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server closed unexpected", err)
		}
	}()

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
	<-quitChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := service.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

type server struct {
	gin *gin.Engine

	Addr string

	isInit bool
}

func Server() *server {
	return s
}

func (s *server) Engine() *gin.Engine {
	if !s.isInit {
		return nil
	}
	return s.gin
}

func (s *server) addDefaultRoute() {
	s.gin.GET("/ping", ping)
}

// Ping handler
func ping(ctx *gin.Context) {
	response.Success(ctx, "ok")
}
