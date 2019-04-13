package plugin

import "MCDaemon-go/plugins/BackupPlugin"

func CreatePluginsList(isload bool) (PluginMap, PluginMap) {
	//可使用插件列表
	PluginsList := make(PluginMap)
	//已被禁用插件列表
	DisablePluginsList := make(PluginMap)

	//加载热插件
	PluginsList.GetHotPlugins(isload)

	//注册冷插件
	PluginsList.RegisterPlugin("!!server", &BasePlugin{})                //基础插件
	PluginsList.RegisterPlugin("!!backup", &BackupPlugin.BackupPlugin{}) //备份插件插件
	PluginsList.RegisterPlugin("!!yinyinmaster", &Yinyin{})              //例子插件
	PluginsList.RegisterPlugin("!!image", &ImagePlugin{})                //镜像插件
	PluginsList.RegisterPlugin("!!SDChat", &SDChatPlugin{})              //沙雕聊天机器人插件
	PluginsList.RegisterPlugin("!!tps", &TpsPlugin{})                    //tps插件

	return PluginsList, DisablePluginsList
}
