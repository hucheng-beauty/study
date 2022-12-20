package main

import (
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

/*

	统一的错误处理逻辑
*/

/*
	http服务器的性能分析
		import _ "net/http/pprof"
		访问/debug/pprof
		使用go tool pprof
		使用 go tool pprof 分析性能
*/

type userErr string

func (u userErr) Error() string {
	return u.Message()
}

func (u userErr) Message() string {
	return string(u)
}

func HandleFileList(writer http.ResponseWriter, request *http.Request) error {
	path := request.URL.Path[len("/list/"):]
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

type appHandler func(writer http.ResponseWriter, request *http.Request) error

// 函数式编程的应用: errWrapper
func errWrapper(handler appHandler) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("[errWrapper] panic: ", r)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		if err := handler(writer, request); err != nil {
			log.Println("[errWrapper] error: ", err.Error())

			if userError, ok := err.(userErr); ok {
				http.Error(writer, userError.Message(), http.StatusBadRequest)
				return
			}

			code := http.StatusInternalServerError
			switch {
			case os.IsNotExist(err):
				code = http.StatusNotFound
			case os.IsPermission(err):
				code = http.StatusForbidden
			default:
			}

			http.Error(writer, http.StatusText(code), code)
		}
	}
}

func main() {
	// 文件服务器列表
	http.HandleFunc("/list/", errWrapper(HandleFileList))

	if err := http.ListenAndServe(":8888", nil); err != nil {
		panic(err)
	}
}
