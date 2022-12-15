package main

import (
	"flag"

	"github.com/AllPaste/web-bbf/config"
	"github.com/AllPaste/web-bbf/internal/server"
	_ "go.uber.org/automaxprocs"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "path", "./config/config.yaml", "config path file.")
}

func main() {
	// 加载配置
	flag.Parse()
	config.LoadConfig(configPath)

	// 启动服务
	server.Run(&config.Cfg)
}
