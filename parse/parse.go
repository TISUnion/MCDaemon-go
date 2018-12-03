/**
命令语法解析器
*/

package parse

import (
	"regexp"
	"strings"
)

type Command struct {
	Player string   //发出命令的玩家
	Cmd    string   //命令
	Argv   []string //参数
	Type   int      //命令类型 系统级还是用户级; 0代表系统级, 1代表用户级
}

//解析玩家输入文字，判断是否是命令
func (c *Command) Parse(world string) bool {
	re := regexp.MustCompile(`\[\d+:\d+:\d+\]\s+\[Server thread/INFO\]:\s+<(?P<player>.+)>\s+(?P<commond>((!|!!)+.+))\s*`)
	match := re.FindStringSubmatch(world)
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
		if strings.Contains(c.Cmd, "!") {
			c.Type = 1
		} else {
			c.Type = 0
		}
		return true
	}
	return false
}
