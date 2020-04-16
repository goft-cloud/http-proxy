package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goft-cloud/http-proxy/config"
	"github.com/goft-cloud/http-proxy/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	s *server
)

type server struct {
	gin *gin.Engine

	Addr string

	isInit bool
}

func Server() *server {
	return s
}

// Engine 返回engine对象
func (s *server) Engine() *gin.Engine {
	if !s.isInit {
		return nil
	}
	return s.gin
}

// 初始化 Server
func Init() {
	s = &server{
		gin:    gin.New(),
		isInit: true,
	}

	// 默认路由
	s.addDefaultRoute()

	// server 地址
	s.Addr = fmt.Sprintf("%s:%d", config.App.Host, config.App.Port)

}

// 启动 server
func Run() error {

	// 创建server
	return newServer()
}

// 创建服务
func newServer() error {
	s.gin.Use(middleware.Request, middleware.Log, middleware.Recovery)

	service := &http.Server{
		Addr:    s.Addr,
		Handler: s.gin,
	}

	fmt.Println(s.Addr)

	go func() {
		if err := service.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("server closed unexpected", err)
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

func (s *server) addDefaultRoute() {
	s.gin.GET("/ping", ping)
}

// ping ping handler
func ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "ping",
	})
}
