package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func httpPost() {
	resp, err := http.Post("http://www.01happy.com/demo/accept.php",
		"application/x-www-form-urlencoded",
		strings.NewReader("name=cjb"))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func main() {
	body := struct {
		method string
		id     string
	}{
		method: "HelloService",
		id:     "0",
	}
	marshal, err := json.Marshal(body)
	if err != nil {
		return
	}

	resp, err := http.Post("http://localhost:1234/jsonrpc", "application/json", strings.NewReader(string(marshal)))
	if err != nil {
		return
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(bytes))

}
