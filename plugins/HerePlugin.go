/*
 * 统计信息助手
 * author: Sciroccogti
 * 参考：https://github.com/TISUnion/Here
 */

package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/lib"
)

type HerePlugin struct {
	highlight   bool              // 是否开启高亮功能
	player      string            // 使用插件的玩家
	dim_convert map[string]string // 编号转维度
}

func (p *HerePlugin) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) == 0 {
		if p.highlight {
			s.Tell(c.Player, command.Text{"您将会被高亮15秒", "yellow"})
			s.Execute("/effect give " + c.Player + " minecraft:glowing 15 1 true")
		}
		s.Execute("/data get entity " + c.Player)
		p.player = c.Player

	} else {
		switch c.Argv[0] {
		case "res":
			if len(c.Argv) == 5 { // && c.Argv[1] == p.player {
				posstr := "[x:" + c.Argv[2] + ", y:" + c.Argv[3] + ", z:" + c.Argv[4] + "]"
				s.Say(command.Text{p.player + "在" + p.dim_convert[c.Argv[1]] + posstr + "向你问好", "green"})
				p.player = ""
			}

		case "set":
			if len(c.Argv) < 2 {
				s.Tell(c.Player, command.Text{"请输入 on 或 off！", "red"})
			} else {
				if c.Argv[1] == "off" {
					p.highlight = false
					s.Say(command.Text{"插件 here 高亮功能已关闭", "yellow"})
				} else if c.Argv[1] == "on" {
					p.highlight = true
					s.Say(command.Text{"插件 here 高亮功能已开启", "yellow"})
				} else {
					s.Tell(c.Player, command.Text{"请输入 on 或 off！", "red"})
				}
			}

		default:
			here1 := command.Text{"!!here", "white"}
			here2 := command.Text{"广播自己的位置\n", "green"}
			set1 := command.Text{"!!here set [on/off]", "white"}
			set2 := command.Text{"开/关高亮功能，默认开\n", "green"}
			s.Tell(c.Player, here1, here2, set1, set2)
		}
	}

}

func (p *HerePlugin) Init(s lib.Server) {
	p.highlight = true
	p.dim_convert = make(map[string]string)
	p.dim_convert["0"] = "主世界"
	p.dim_convert["1"] = "末地"
	p.dim_convert["-1"] = "地狱"
}

func (p *HerePlugin) Close() {
}
