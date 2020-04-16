package proxy

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func DoProxy(c *gin.Context) {

	request := c.Request
	writer := c.Writer

	start := time.Now().UnixNano() / 1000000

	url := request.Header.Get("target-addr")

	fmt.Printf("url is=" + url + "\n")

	req, _ := http.NewRequest(request.Method, url, request.Body)

	//rs,_:=ioutil.ReadAll(request.Body)
	//
	//fmt.Printf("rbody=%s \n",rs)

	for k, v := range request.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	resp, err := client.Do(req)
	if (err != nil) {
		fmt.Printf("error=" + err.Error() + " \n")
		return
	}

	defer resp.Body.Close()

	for k1, v1 := range resp.Header {
		for _, vv1 := range v1 {
			writer.Header().Add(k1, vv1)
		}
	}

	data, _ := ioutil.ReadAll(resp.Body)

	end := time.Now().UnixNano() / 1000000

	fmt.Printf(" cost=" + strconv.Itoa(int(end-start)) + "\n")

	//fmt.Println("response=" + string(data))
	_, _ = writer.Write(data)
}
