package server

import (
	"MCDaemon-go/command"
	"MCDaemon-go/config"
	"MCDaemon-go/lib"
	"fmt"
	"strings"
)

//运行所有语法解析器
func (svr *Server) RunParsers(word string) {
	for _, val := range svr.parserList {
		cmd, ok := val.Parsing(word)
		if ok && svr.pluginList[cmd.Cmd] != nil {
			//异步运行插件
			svr.pulginPool <- 1
			if cmd.Player != "" {
				svr.WriteLog("info", fmt.Sprintf("玩家 %s 运行了 %s %s命令", cmd.Player, cmd.Cmd, strings.Join(cmd.Argv, " ")))
			}
			go svr.RunPlugin(cmd)
		}
	}
}

//运行插件
func (svr *Server) RunPlugin(cmd *command.Command) {
	svr.pluginList[cmd.Cmd].Handle(cmd, svr)
	<-svr.pulginPool
}

//等待现有插件的完成并停止后面插件的运行，在执行相关操作
func (svr *Server) RunUniquePlugin(handle func()) {
	<-svr.pulginPool
	//根据插件最大并发数进行堵塞
	maxRunPlugins, _ := config.Cfg.Section("MCDeamon").Key("maxRunPlugins").Int()
	for i := 0; i < maxRunPlugins; i++ {
		svr.pulginPool <- 1
	}
	handle()
	for i := 0; i < maxRunPlugins; i++ {
		<-svr.pulginPool
	}
	svr.pulginPool <- 1
}

//获取当前实例的插件列表
func (svr *Server) GetPluginList() map[string]lib.Plugin {
	return svr.pluginList
}

//获取当前实例的禁用插件列表
func (svr *Server) GetDisablePluginList() map[string]lib.Plugin {
	return svr.disablePluginList
}

//获取语法解析器列表
func (svr *Server) GetParserList() []lib.Parser {
	return svr.parserList
}
