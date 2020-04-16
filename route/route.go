package route

import (
	"github.com/goft-cloud/http-proxy/proxy"
	"github.com/goft-cloud/http-proxy/server"
)

func Init() error {
	router := server.Server().Engine()

	router.Any("/proxy", proxy.DoProxy)
	return nil
}
