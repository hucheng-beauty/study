package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// http request body only read once
func readBodyOnce(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read form body failed:%v", err)
		return
	}
	fmt.Fprintf(w, "read the data %s\n", string(body))
}

func getBodyIsNil() {

}

func main() {
	http.HandleFunc("/readBodyOnce", readBodyOnce)
	http.ListenAndServe("8080", nil)
}
