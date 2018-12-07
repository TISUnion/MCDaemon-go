package plugin

var PluginsList PluginMap

func init() {
	//加载热插件
	PluginsList.GetHotPlugins(false)

	//注册冷插件
	PluginsList.RegisterPlugin("!!yinyinmaster", Yinyin{}) //测试插件
}
