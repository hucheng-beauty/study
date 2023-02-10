package main

import (
	"context"
	"fmt"
	"net/http"
)

func serveApp(addr string, handel http.Handler, stop chan struct{}) error {
	s := http.Server{Addr: addr, Handler: handel}

	go func() {
		<-stop
		s.Shutdown(context.Background())
	}()
	return s.ListenAndServe()
}

func serveDebug(addr string, handel http.Handler, stop chan struct{}) error {
	s := http.Server{Addr: addr, Handler: handel}

	go func() {
		<-stop // wait for stop signal
		s.Shutdown(context.Background())
	}()
	return s.ListenAndServe()
}

// serve just apply one service
func serve(addr string, handel http.Handler, stop chan struct{}) error {
	s := http.Server{Addr: addr, Handler: handel}

	go func() {
		<-stop
		s.Shutdown(context.Background())
	}()

	return s.ListenAndServe()
}

func app() {
	quit := make(chan error, 2)
	stop := make(chan struct{})

	go func() {
		quit <- serveApp(":8080", nil, stop)
	}()

	go func() {
		quit <- serveDebug(":8081", nil, stop)
	}()

	var stopped bool
	for i := 0; i < cap(quit); i++ {
		if err := <-quit; err != nil {
			fmt.Printf("error: %v\n", err)
		}
		if !stopped {
			stopped = true
			close(stop)
		}
	}
}
