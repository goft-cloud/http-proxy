package main

import (
	"fmt"
	"github.com/goft-cloud/http-proxy/config"
	"github.com/goft-cloud/http-proxy/log"
	"github.com/goft-cloud/http-proxy/route"
	"github.com/goft-cloud/http-proxy/server"
	"os"
)

func main() {
	if err := server.Run(); err != nil {
		fmt.Println("Start Application fail!" + err.Error())
	}
}

func init() {
	if err := bootstrap(); err != nil {
		fmt.Println(" Initialize Application Fail!" + err.Error())
		os.Exit(1)
	}

	addr := server.Server().Addr
	fmt.Println("Start Application Success")
	fmt.Println("Listen:" + addr)
}

func bootstrap() error {
	fmt.Print("Initialize config ...")
	if err := config.Init(); err != nil {
		return err
	}
	fmt.Println("Initialize config success!")

	fmt.Println("Initialize log ...")
	if err := log.Init(); err != nil {
		return err
	}
	fmt.Println("Initialize log success!")

	fmt.Println("Initialize server ...")
	if err := server.Init(); err != nil {
		return err
	}
	fmt.Println("Initialize server success!")

	fmt.Println("Initialize route ...")
	if err := route.Init(); err != nil {
		return err
	}
	fmt.Println("Initialize route success!")

	return nil
}
