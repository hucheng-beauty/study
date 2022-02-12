package main

import (
	"log"
	"net/http"
)

/*
	Body: 只能读一次
*/

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//bytes1, err := io.ReadAll(r.Body)
		//if err != nil {
		//	return
		//}
		//fmt.Fprintf(w, string(bytes1))
		//
		//bytes2, err := io.ReadAll(r.Body)
		//if err != nil {
		//	fmt.Fprintf(w, "second error:")
		//}
		//fmt.Fprintf(w, string(bytes2))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
