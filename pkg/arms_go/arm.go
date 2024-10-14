package main

//
// import (
//     "fmt"
//     "net/http"
//
//     arms "github.com/aliyun/arms-go-agent"
// )
//
// func main() {
//     // 初始化 ARMS Go Agent
//     config := arms.Config{
//         AppName:    "your-app-name",
//         LicenseKey: "your-license-key",
//     }
//     err := arms.InitARMS(&config)
//     if err != nil {
//         fmt.Printf("Failed to initialize ARMS Go Agent: %v\n", err)
//         return
//     }
//     defer arms.Shutdown()
//
//     // 创建一个简单的 HTTP 服务器
//     http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//         // 使用 ARMS Go Agent 创建一个新的 span
//         ctx := r.Context()
//         span := arms.StartSpan(ctx, "handle-request")
//         defer span.End()
//
//         // 添加一些自定义属性
//         span.SetAttribute("path", r.URL.Path)
//
//         // 处理请求
//         w.Write([]byte("Hello, ARMS Go Agent!"))
//     })
//
//     // 启动 HTTP 服务器
//     fmt.Println("Starting server on :8080")
//     if err := http.ListenAndServe(":8080", nil); err != nil {
//         fmt.Printf("ListenAndServe error: %v\n", err)
//     }
// }
