package lib

type Command struct {
	Player     string   //玩家名
	PluginName string   //插件名
	Argv       []string //参数
	Cmd        string   //命令名
}

func (c *Command) GetPluginName() {

}
