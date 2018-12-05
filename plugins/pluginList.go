package plugin

import (
	"MCDaemon-go/config"
	"MCDaemon-go/lib"
)

var PluginsList map[string]lib.Plugin

func init() {
	Plugins := config.GetPlugins(false)
	for k, v := range Plugins {
		PluginsList[k] = HotPlugin(v)
	}
}
