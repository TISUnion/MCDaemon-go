package main

import (
	"MCDaemon-go/config"
	"MCDaemon-go/server"
)

var (
	MCDeamon []string
)

func init() {
	//配置eula文件
	config.SetEula()
	//获取所有启动项配置
	MCDeamon = config.GetStartConfig()
}

func main() {
	//获取服务器实例
	svr := server.GetServerInstance()
	//初始化
	svr.Init(MCDeamon)
	//等待加载地图
	svr.WaitEndLoading()
	//正式运行MCD
	svr.Run()
	defer svr.Close()
}
