package plugin

var PluginsList PluginMap

func init() {

	PluginsList = make(PluginMap)

	//加载热插件
	PluginsList.GetHotPlugins(false)

	//注册冷插件
	PluginsList.RegisterPlugin("!!server", BasePlugin{})   //基础插件
	PluginsList.RegisterPlugin("!!yinyinmaster", Yinyin{}) //例子插件
}
