package main

import (
	"MCDaemon-go/config"
	"MCDaemon-go/container"
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
	c := container.GetInstance()
	defaultServer := &server.Server{}
	//初始化
	defaultServer.Init("default", MCDeamon)
	//加入到容器中
	c.Add("default", defaultServer)
	//等待加载地图
	defaultServer.WaitEndLoading()
	//正式运行MCD
	defaultServer.Run()

}
