package plugin

//冷加载命令列表
var Commands map[string]string

func init() {
	Commands = make(map[string]string)
	Commands = map[string]string{}
}
