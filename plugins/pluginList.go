package plugin

func CreatePluginsList() (PluginMap, PluginMap) {
	//可使用插件列表
	PluginsList := make(PluginMap)
	//已被禁用插件列表
	DisablePluginsList := make(PluginMap)

	//加载热插件
	PluginsList.GetHotPlugins(false)

	//注册冷插件
	PluginsList.RegisterPlugin("!!server", &BasePlugin{})   //基础插件
	PluginsList.RegisterPlugin("!!yinyinmaster", &Yinyin{}) //例子插件

	return PluginsList, DisablePluginsList
}
