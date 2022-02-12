package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

/*
	http:
		使用 http 客户端发送请求
		使用 http.Client 控制请求头部等
		使用 httputil简化工作
*/

const url = "http://www.imooc.com"

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"

func main() {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	request.Header.Add("UserAgent", userAgent)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	s, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", s)
}
