package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/goft-cloud/http-proxy/config"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	client      *http.Client
	once        sync.Once
	proxyConfig *PoolConfig
)

type PoolConfig struct {
	Timeout               int `toml:"timeout"`
	KeepAlive             int `toml:"keep_alive"`
	MaxIdleConns          int `toml:"max_idle_conns"`
	IdleConnTimeout       int `toml:"idle_conn_timeout"`
	TLSHandshakeTimeout   int `toml:"tls_handshake_timeout"`
	ExpectContinueTimeout int `toml:"expect_continue_timeout"`
	MaxConnsPerHost       int `toml:"max_conns_per_host"`
	MaxIdleConnsPerHost   int `toml:"max_idle_conns_per_host"`
}

func HttpClient() *http.Client {
	once.Do(func() {
		// 初始化
		initPoolConfig()

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

	fmt.Println(proxyConfig)

	return client
}

func initPoolConfig() {
	proxyConfig = &PoolConfig{
		Timeout:               1,
		KeepAlive:             1,
		MaxIdleConns:          100,
		IdleConnTimeout:       1,
		TLSHandshakeTimeout:   1,
		ExpectContinueTimeout: 1,
		MaxConnsPerHost:       20,
		MaxIdleConnsPerHost:   20,
	}

	if err := config.DecodeKey("proxy", proxyConfig); err != nil {
		fmt.Println("Decode pool config key error1" + err.Error())
	}
}
