FROM golang:1.14-alpine as builder
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct
ENV GOCACHE=/go/pkg/go-build
ARG TZ=Asia/Shanghai
ENV TZ ${TZ}
# 设置字符集
ENV LANG C.UTF-8
# 更新安装源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
# 按照必要软件
RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev
RUN apk update && apk add --no-cache bash curl strace gdb
# 设置项目目录
ARG PROJECT_PATH=/go/src/github.com/king19800105/go-ws-demo
# 设置工作目录
WORKDIR ${PROJECT_PATH}
# 拷贝go module相关文件
COPY go.mod .
COPY go.sum .
# 运行下载go module中的软件
RUN go mod download
# 将服务器的go工程代码加入到docker容器中
ADD . ${PROJECT_PATH}
# 运行编译
RUN go build -a -o hardware cmd/instrument/main.go
# 开启端口
EXPOSE 9601
EXPOSE 9602
# 启动指令
ENTRYPOINT ["./hardware >> /var/log/go/print.log"]