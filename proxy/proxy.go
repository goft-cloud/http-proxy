package proxy

import (
	"github.com/gin-gonic/gin"
	"github.com/goft-cloud/http-proxy/log"
	"github.com/goft-cloud/http-proxy/response"
	"io/ioutil"
	"net/http"
)

const (
	targetHeaderName = "target-addr"
)

// Do proxy
func DoProxy(c *gin.Context) {
	request := c.Request
	writer := c.Writer

	// Target url from header
	url := request.Header.Get(targetHeaderName)
	if len(url) == 0 {
		response.Fatal(c, "Header for `target-addr` is not exist or empty!")
		return
	}

	// New target request
	tRequest, err := http.NewRequest(request.Method, url, request.Body)
	if err != nil {
		response.Fatal(c, "New target request error!"+err.Error())
		return
	}

	// Set target request header
	for tk, tv := range request.Header {
		for _, tvHeader := range tv {
			tRequest.Header.Add(tk, tvHeader)
		}
	}

	// Do target request
	tResponse, err := HttpClient().Do(tRequest)
	if err != nil {
		response.Fatal(c, "New request do error!"+err.Error())
		return
	}

	// Close target response
	defer func() {
		err = tResponse.Body.Close()
		if err != nil {
			log.Error(c, "Close target response error!"+err.Error())
		}
	}()

	// Set target response header
	for rk, rv := range tResponse.Header {
		for _, rvHeader := range rv {
			writer.Header().Add(rk, rvHeader)
		}
	}

	// Target response body
	body, err := ioutil.ReadAll(tResponse.Body)
	if err != nil {
		response.Fatal(c, "Read target body error!"+err.Error())
		return
	}

	_, err = writer.Write(body)
	if err != nil {
		response.Fatal(c, "Write body to target response error!"+err.Error())
	}
}
