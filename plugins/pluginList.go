package plugin

import (
	"MCDaemon-go/config"
	"MCDaemon-go/lib"
)

var PluginsList map[string]lib.Plugin

func init() {
	Plugins := config.GetPlugins(false)
	//加载热插件
	for k, v := range Plugins {
		PluginsList[k] = HotPlugin(v)
	}
}
