/**
*基础插件
*提供服务器停止，启动和重启功能
 */

package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/lib"
	"fmt"
)

type BasePlugin struct {
}

func (hp *BasePlugin) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) == 0 {
		c.Argv = append(c.Argv, "help")
	}
	switch c.Argv[0] {
	case "restart":
		f := s.Restart
		s.RunUniquePlugin(f)
	case "stop":
		fmt.Println("结束")
		f := s.Close
		s.RunUniquePlugin(f)
	case "reload":
		f := func() {
			PluginMap(s.GetPluginList()).GetHotPlugins(true)
		}
		s.RunUniquePlugin(f)
	case "ban":
		if len(c.Argv) > 1 {
			if plugin, ok := PluginMap(s.GetPluginList()).DelPlugin(c.Argv[1]); ok {
				PluginMap(s.GetDisablePluginList()).RegisterPlugin(c.Argv[1], plugin)
			} else {
				s.Tell(c.Player, "不存在该插件，建议!!server show list查看可使用的插件")
			}
		} else {
			s.Tell(c.Player, "请输入插件名称")
		}
	case "pardon":
		if len(c.Argv) > 1 {
			if plugin, ok := PluginMap(s.GetDisablePluginList()).DelPlugin(c.Argv[1]); ok {
				PluginMap(s.GetPluginList()).RegisterPlugin(c.Argv[1], plugin)
			} else {
				s.Tell(c.Player, "不存在该插件，建议!!server show banlist查看已被禁用的插件")
			}
		} else {
			s.Tell(c.Player, "请输入插件名称")
		}
	case "show":
		if len(c.Argv) > 1 {
			if c.Argv[1] == "list" {
				var text string
				for k, _ := range PluginMap(s.GetPluginList()) {
					text += k + "\\n"
				}
				s.Tell(c.Player, "插件列表：\\n"+text)
			} else if c.Argv[1] == "banlist" {
				var text string
				for k, _ := range PluginMap(s.GetDisablePluginList()) {
					text += k + "\\n"
				}
				s.Tell(c.Player, "已禁用插件列表：\\n"+text)
			}
		} else {
			s.Tell(c.Player, "请输入查看插件类型")
		}
	default:
		text := "!!server restart 重启服务器\\n!!server stop 关闭服务器\\n!!server reload 重新加载热插件\\n!!server ban [插件名] 禁用插件\\n!!server pardon [插件名] 恢复禁用的插件\\n!!server show list 查看插件列表\\n!!server show banlist 查看禁用插件列表"
		s.Tell(c.Player, text)
	}
}
