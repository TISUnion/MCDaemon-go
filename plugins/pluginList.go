package plugin

import "MCDaemon-go/config"

//所有加载插件列表
var PluginList map[string]Plugin

func init() {
	PluginList = make(map[string]Plugin)
	PluginList = map[string]Plugin{}
	cg := config.GetInstance()
	PluginFromConfList := cg.GetPlugins()
	for k, v := range PluginFromConfList {
		PluginList[k] = v
	}
}
