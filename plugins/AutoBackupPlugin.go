/*
 * 自动定时备份插件
 * 定时运行：rsync -a --delete minecraft/ back-up/auto
 * 他会对所有人说：嘤嘤嘤
 */

package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/lib"
	"fmt"
)

type AutoBackup struct {
}

func (hp *AutoBackup) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) == 0 {
		c.Argv = append(c.Argv, "help")
	}

	switch c.Argv[0] {
	case "set":
		if len(c.Argv) < 2 {
			s.Tell(c.Player, command.Text{"请输入要设定的小时数！", "red"})
		} else {
			s.Tell(c.Player, command.Text{fmt.Sprintf("收到参数%d", c.Argv[1]), "green"})
		}
		s.Tell(c.Player, command.Text{"本插件还未完工！", "red"})

	case "query":
		s.Tell(c.Player, command.Text{"本插件还未完工！", "red"})

	default:
		set1 := command.Text{"!!autobk set [小时数]", "white"}
		set2 := command.Text{"设定存档间隔\\n", "green"}
		query1 := command.Text{"!!autobk query", "white"}
		query2 := command.Text{"查询上次存档时间\\n", "green"}
		s.Tell(c.Player, set1, set2, query1, query2)
	}
}

func (hp *AutoBackup) Init(s lib.Server) {
}

func (hp *AutoBackup) Close() {
}
