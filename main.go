package main

import (
	"github.com/goft-cloud/http-proxy/config"
	"github.com/goft-cloud/http-proxy/log"
	"github.com/goft-cloud/http-proxy/route"
	"github.com/goft-cloud/http-proxy/server"
)

func main() {
	// 启动初始化
	config.Init()

	log.Init()

	server.Init()

	route.Init()

	server.Run()
}
