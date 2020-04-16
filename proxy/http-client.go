package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"
)

var client *http.Client
var once sync.Once

func HttpClient() *http.Client {
	once.Do(func() {
		netDialer := &net.Dialer{
			Timeout:   1 * time.Second,
			KeepAlive: 1 * time.Second,
		}

		transport := &http.Transport{
			//Proxy:                 proxy,
			DialContext:           netDialer.DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       1 * time.Second,
			TLSHandshakeTimeout:   1 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxConnsPerHost:       20,
			MaxIdleConnsPerHost:   20,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		}

		client = &http.Client{
			Transport: transport,
		}
	})

	return client
}
