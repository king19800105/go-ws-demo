package config

import (
	"flag"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go/build"
	"os"
	"strings"
)

const msgConfig = "配置对象初始化错误"

// 初始化viper配置
func Viperize() (v *viper.Viper, err error) {
	v = viper.New()
	if err = bindFlags(v); nil != err {
		err = errors.Wrap(err, msgConfig)
		return
	}
	configureViper(v)
	return
}

// 读取指定配置文件，并合并
func LoadFile(v *viper.Viper, path string, env string) error {
	basePath := os.Getenv("GOPATH")
	if "" == basePath {
		basePath = build.Default.GOPATH
	}

	fileName := "config." + env
	v.SetConfigName(fileName)
	v.AddConfigPath(basePath + path)
	return v.ReadInConfig()
}

func bindFlags(v *viper.Viper) error {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := v.BindPFlags(pflag.CommandLine); nil != err {
		return errors.Wrap(err, msgConfig)
	}

	return nil
}

func configureViper(v *viper.Viper) {
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
}
