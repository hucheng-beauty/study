package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

/*
	goroutine 注意事项:
		1.Keep yourself busy or do the work yourself.
		2.Leave concurrency to the caller.
		3.Never start a goroutine without knowing when it will stop.
*/

func serve(addr string, handler http.Handler, stop <-chan struct{}) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		<-stop // wait for stop signal
		s.Shutdown(context.Background())
	}()

	return s.ListenAndServe()
}

func serveApp(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Hello, QCon!")
	})
	return serve("0.0.0.0:8080", mux, stop)
}

func serveDebug(stop <-chan struct{}) error {
	return serve("127.0.0.1:8001", http.DefaultServeMux, stop)
}

func main() {
	/*
		// no graceful style
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, GopherCon SG")
		})
		go func() {
			if err := http.ListenAndServe(":8080", nil); err != nil {
				log.Fatal(err)
			}
		}()

		select {}
	*/

	// if your goroutine cannot make progress until it gets the result from another,
	// oftentimes it is simpler to just do the work yourself rather than to delegate it.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, GopherCon SG")
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	/*
		// return the content of directory

		// question: 同步调用,阻塞直到所有目录条目都被读取;可能需要更长的时间和内存来构建切片
		func ListDirectory(dir string) ([]string, error)

		// question:
		func ListDirectory(dir string) chan string

		// Leave concurrency to the caller.
		func ListDirectory(dir string, fn func(string))
	*/

	// Never start a goroutine without knowing when it will stop.
	done := make(chan error, 2)
	stop := make(chan struct{})
	go func() {
		done <- serveDebug(stop)
	}()
	go func() {
		done <- serveApp(stop)
	}()

	var stopped bool
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			fmt.Printf("error: %v\n", err)
		}
		if !stopped {
			stopped = true
			close(stop)
		}
	}
}
