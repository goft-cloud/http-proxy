package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

var client *http.Client
var transport *http.Transport

func init() {
	//proxy := func(request *http.Request) (*url.URL, error) {
	//
	//	forwardAddr := request.Header.Get("X-FORWARD-FOR")
	//
	//	fmt.Println(" forwardAddr="+forwardAddr)
	//	u, _ := url.Parse(forwardAddr)
	//	fmt.Println("--" + u.Scheme + "---" + u.Host + "---" + u.Port())
	//
	//	addr, _ := net.ResolveIPAddr("ip", u.Host)
	//
	//	fmt.Println(" add=" + addr.String() + "\n")
	//
	//	//return url.Parse("http://127.0.0.1:8080")
	//	//return url.Parse("http://47.108.140.202")
	//	return url.Parse(u.Scheme + "://" + addr.String())
	//	//return url.Parse("https://14.215.140.116:443")
	//}

	netDialer := &net.Dialer{
		Timeout:   1 * time.Second,
		KeepAlive: 1 * time.Second,
	}

	transport = &http.Transport{
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
}