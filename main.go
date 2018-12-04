package main

import (
	"MCDaemon-go/config"
	"MCDaemon-go/server"
	"fmt"
)

var (
	commandArgv []string //服务器启动参数
	svr         *server.Server
	conf        *config.Config
)

func init() {
	//判断eula是否为真
	config.SetEula()
	//获取conig实例
	conf = config.GetInstance()
	//加载服务器启动配置
	MCDconfig := conf.GetConfig()
	commandArgv = []string{
		MCDconfig["Xmx"],
		MCDconfig["Xms"],
		"-jar",
		fmt.Sprintf("%s/%s", MCDconfig["serverPath"], MCDconfig["serverName"]),
	}
	if MCDconfig["gui"] != "true" {
		commandArgv = append(commandArgv, "nogui")
	}
}

func main() {
	// 创建服务器实例并启动
	svr = server.GetServerInstance()
	//初始化服务器
	svr.Init(commandArgv)
	// 等待地图加载
	svr.WaitEndLoading()
	// 运行MCD
	svr.Run()
	defer svr.Close()
}
