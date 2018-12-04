package command

import "MCDaemon-go/config"

type Command struct {
	Player     string   //玩家名
	PluginName string   //插件名
	Argv       []string //参数名
	Cmd        string   //命令名
}

func (c *Command) GetPluginName() {
	conf := config.GetInstance()
	plugins := conf.GetPlugins()
	c.PluginName = plugins[c.Cmd]
}
