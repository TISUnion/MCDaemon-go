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
	case "help":
		text := "!!server restart 重启服务器\\n!!server stop 关闭服务器"
		s.Tell(text, c.Player)
	case "restart":
		f := s.Restart
		s.RunUniquePlugin(f)
	case "stop":
		fmt.Println("结束")
		s.Close()
	}
}
