# 使用 golang 官方提供的镜像作为基础镜像
FROM golang:1.16.4

# 设置工作目录
WORKDIR /app

# 将应用程序代码复制到镜像中
COPY main.go .

# 编译应用程序并生成可执行文件
RUN go build -o myapp main.go

# 设置容器启动时需要执行的命令
CMD ["./myapp"]
