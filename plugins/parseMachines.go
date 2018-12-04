package plugin

import "MCDaemon-go/parse"

var ParseMachines []parse.ParseMachine

//冷加载插件列表
func init() {
	ParseMachines = []parse.ParseMachine{}
}
