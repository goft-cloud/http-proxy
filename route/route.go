package route

import (
	"github.com/goft-cloud/http-proxy/proxy"
	"github.com/goft-cloud/http-proxy/server"
)

func Init() error {
	router := server.Server().Engine()

	// Add proxy uri
	router.Any("/", proxy.DoProxy)
	return nil
}
