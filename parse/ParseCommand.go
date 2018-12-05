package parse

import (
	"MCDaemon-go/lib"
	"regexp"
	"strings"
)

//解析玩家输入文字，判断是否是命令 ， 实现了ParseMachine接口
type ParseCommand struct {
	Player string   //发出命令的玩家
	Cmd    string   //命令
	Argv   []string //参数
}

//实现解析器接口
func (c ParseCommand) Parsing(word string) (*lib.Command, bool) {
	re := regexp.MustCompile(`\[\d+:\d+:\d+\]\s+\[Server thread/INFO\]:\s+<(?P<player>.+)>\s+(?P<commond>((!|!!)+.+))\s*`)
	match := re.FindStringSubmatch(word)
	groupNames := re.SubexpNames()

	result := make(map[string]string)

	//匹配到命令时
	if len(match) != 0 {
		// 转换为map
		for i, name := range groupNames {
			result[name] = match[i]
		}

		// 解析命令以及参数
		cmdArgv := strings.Fields(result["commond"])
		c.Player = result["player"]
		c.Cmd = cmdArgv[0]
		c.Argv = cmdArgv[1:]
		_commond := &lib.Command{
			Player: c.Player,
			Cmd:    c.Cmd,
			Argv:   c.Argv,
		}
		_commond.GetPluginName()
		return _commond, true
	}
	return nil, false
}
