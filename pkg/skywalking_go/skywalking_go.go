package main

//
// import (
//     "log"
//     "net/http"
//
//     "github.com/SkyAPM/go2sky"
//     httpPlugin "github.com/SkyAPM/go2sky/plugins/http"
//     "github.com/SkyAPM/go2sky/reporter"
// )
//
// func main() {
//     // Create a new reporter
//     r, err := reporter.NewGRPCReporter("localhost:11800")
//     if err != nil {
//         log.Fatalf("new reporter error %v \n", err)
//     }
//     defer r.Close()
//
//     // Create a new tracer
//     tracer, err := go2sky.NewTracer("skywalking-go-service", go2sky.WithReporter(r))
//     if err != nil {
//         log.Fatalf("create tracer error %v \n", err)
//     }
//
//     // Create an HTTP server with SkyWalking instrumentation
//     mux := http.NewServeMux()
//     mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//         w.Write([]byte("Hello, SkyWalking!"))
//     })
//
//     httpHandler, err := httpPlugin.NewHandler(mux, tracer)
//     if err != nil {
//         log.Fatalf("create http handler error %v \n", err)
//     }
//
//     // Start the HTTP server
//     log.Println("Starting server on :8080")
//     err = http.ListenAndServe(":8080", httpHandler)
//     if err != nil {
//         log.Fatalf("ListenAndServe error %v \n", err)
//     }
// }
