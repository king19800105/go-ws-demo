package main

import (
	"fmt"
	"github.com/king19800105/go-ws-demo/pkg/config"
	"github.com/king19800105/go-ws-demo/pkg/hardware/instrument"
	"github.com/king19800105/go-ws-demo/pkg/server"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	msgConfig    = "配置文件初始化失败"
	msgServer    = "插座服务（socket）结束 : %v"
	msgInterrupt = "信号量中断，%s"
)

// 入口
func main() {
	if err := run(); nil != err {
		log.Fatalf(msgServer, err)
	}
}

// 启动服务
func run() (err error) {
	pflag.StringP("http", "h", ":9601", "websocket listen address")
	pflag.StringP("env", "e", "dev", "env set")
	pflag.StringP("cfg-path", "c", "/src/github.com/king19800105/go-ws-demo/configs", "config path")

	// 配置文件
	var cfg *viper.Viper
	{
		cfg, err = createConfig()
		if nil != err {
			return errors.Wrap(err, msgConfig)
		}
	}

	// 设备对象
	var ins = instrument.NewInstrument()
	var errCh = make(chan error)
	// 创建服务
	srv := server.NewServer(cfg.GetString("http"), cfg.GetString("app.route"))
	// 启动服务
	go func(errCh chan error) {
		errCh <- srv.Start(ins)
	}(errCh)
	// 监听中断
	go func(errCh chan error) {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errCh <- fmt.Errorf(msgInterrupt, <-c)
	}(errCh)

	return <-errCh
}

// 初始化配置文件
func createConfig() (cfg *viper.Viper, err error) {
	if cfg, err = config.Viperize(); nil != err {
		return
	}

	path := cfg.GetString("cfg-path")
	env := cfg.GetString("env")
	if err = config.LoadFile(cfg, path, env); nil != err {
		return
	}

	return
}
