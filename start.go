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
	//加入到容器中
	c.Add("default", defaultServer)
	//开启服务器
	c.Group.Add(1)
	go defaultServer.Start("default", MCDeamon, config.Cfg.Section("MCDeamon").Key("server_path").String())
	c.Group.Wait()
}
