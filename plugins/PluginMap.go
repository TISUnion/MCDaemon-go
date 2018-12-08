package plugin

import (
	"MCDaemon-go/config"
	"MCDaemon-go/lib"
)

type PluginMap map[string]lib.Plugin

//加载热插件
func (pm PluginMap) GetHotPlugins(is_reload bool) {
	plugins := config.GetPlugins(is_reload)
	for k, v := range plugins {
		pm[k] = HotPlugin(v)
	}
}

//注册冷插件
func (pm PluginMap) RegisterPlugin(name string, lp lib.Plugin) {
	pm[name] = lp
}
