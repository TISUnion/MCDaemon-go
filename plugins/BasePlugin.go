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

func (hp BasePlugin) Handle(c *command.Command, s lib.Server) {
	switch c.Argv[0] {
	case "restart":
		f := s.Restart
		s.RunUniquePlugin(f)
	case "stop":
		fmt.Println("结束")
		s.Close()
	case "reload":
		f := func() {
			PluginsList.GetHotPlugins(true)
		}
		s.RunUniquePlugin(f)
	case "help":
		text := "!!server restart 重启服务器\\n!!server stop 关闭服务器\\n!!server reload 重新加载热插件"
		s.Tell(text, c.Player)
	default:
		text := "!!server restart 重启服务器\\n!!server stop 关闭服务器\\n!!server reload 重新加载热插件"
		s.Tell(text, c.Player)
	}
}
