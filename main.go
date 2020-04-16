package main

import (
	"fmt"
	"github.com/goft-cloud/http-proxy/config"
	"github.com/goft-cloud/http-proxy/log"
	"github.com/goft-cloud/http-proxy/route"
	"github.com/goft-cloud/http-proxy/server"
)

func init() {
	bootstrap()
}

func main() {
	server.Run()
}

func bootstrap() error {
	fmt.Println("Initialize config ...")
	if err := config.Init(); err != nil {
		return err
	}
	fmt.Println("Initialize config success!")

	log.Init()

	server.Init()

	route.Init()

	return nil
}
