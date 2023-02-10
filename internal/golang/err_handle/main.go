/*
	日志记录与错误无关且对调试没有帮助的信息视为噪音
	只处理错误一次,记录日志做降级处理,或者直接返回错误

	handling error
		如果库有多个调用者,使用原始的错误信息，不要使用 errors.Wrap 保存堆栈信息
		可以处理的错误直接处理,处理不了的错误包装额外的上下文信息往上抛
		错误一旦处理,错误就不在是错误,就不能返回错误值,应该只返回降级数据,然后返回 nil
	github.com/pkg/errors
		通过 errors.New 或者 errors.Errorf 返回错误
		通过 errors.Wrap 或者 errors.Wrapf 保存堆栈信息; 通过 %+v 将堆栈详情记录
		通过 errors.Cause 获取 root error,再进行和 sentinel error 判定
		通过 errors.WithMessage 对错误进行包装
	go1.13
		errors.Is: 等值判定
		errors.As: error 的类型转换
		通过 fmt.Errorf + %w 扩展的错误可以通过 Is 和 As 进行转化或者等值判定

*/
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type userError interface {
	error
	Message() string
}

func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		// 保存堆栈信息
		return nil, errors.Wrap(err, "open failed")
	}
	defer file.Close()

	return []byte("response stream data"), nil
}

func ReadConfig() ([]byte, error) {
	home := os.Getenv("HOME")
	config, err := ReadFile(filepath.Join(home, ".settings.xml"))

	// 对错误进行包装,往上抛
	return config, errors.WithMessage(err, "could not read config")
}

func main() {
	_, err := ReadConfig()
	if err != nil {
		fmt.Printf("original error: %T\n", errors.Cause(err))
		fmt.Printf("stack trace: %+v\n", err)
		os.Exit(1)
	}
}
