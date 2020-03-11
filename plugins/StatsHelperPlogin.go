/*
 * 统计信息助手
 * author: Sciroccogti
 */

package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/lib"
	"fmt"
)

type StatsHelper struct {
	statsmap map[string]string // [代号] 显示名称
}

func (p *StatsHelper) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) == 0 {
		c.Argv = append(c.Argv, "help")
	}

	switch c.Argv[0] {
	case "list":
		s.Tell(c.Player, command.Text{"现有榜单如下（[榜单代号] 显示名称）：", "yellow"})
		for stats := range p.statsmap {
			s.Tell(c.Player, command.Text{fmt.Sprintf("[%s] %s", stats, p.statsmap[stats]), "yellow"})
		}

	case "set":
		if len(c.Argv) < 2 {
			s.Tell(c.Player, command.Text{"请输入榜单代号！", "red"})
		} else {
			if c.Argv[1] == "off" {
				s.Execute(fmt.Sprintf("/scoreboard objectives setdisplay sidebar"))
				s.Say(command.Text{"榜单已关闭显示", "yellow"})
			} else {
				board, ok := p.statsmap[c.Argv[1]]
				if ok {
					s.Execute(fmt.Sprintf("/scoreboard objectives setdisplay sidebar %s", c.Argv[1]))
					s.Say(command.Text{"榜单已显示为 " + board, "yellow"})
				} else {
					s.Tell(c.Player, command.Text{"请检查输入的榜单代号！", "red"})
				}
			}
		}

	case "help":
		list1 := command.Text{"!!stats list", "white"}
		list2 := command.Text{"列出所有榜单\n", "green"}
		set1 := command.Text{"!!stats set [榜单代号(英文)]", "white"}
		set2 := command.Text{"设置要显示的榜单，off为关闭\n", "green"}
		s.Tell(c.Player, list1, list2, set1, set2)
	}

}

func (p *StatsHelper) Init(s lib.Server) {
	p.statsmap = make(map[string]string)
	p.statsmap["deathcount"] = "死亡榜"
	p.statsmap["onlinecount"] = "在线榜"
	p.statsmap["killcount"] = "杀怪榜"
	// TODO: 自动获取服务器中的榜单
}

func (p *StatsHelper) Close() {
}
