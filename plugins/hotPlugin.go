package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/lib"
	"fmt"
	"os/exec"
	"strings"
)

//热加载插件类型
type HotPlugin struct {
	string
}

func (hp *HotPlugin) Handle(c *command.Command, s lib.Server) {
	commandName := "./hotPlugins/" + c.PluginName
	pluginProcess := exec.Command(commandName, c.Argv...)
	buffer, err := pluginProcess.Output()
	if err != nil {
		s.Tell(fmt.Sprint("插件出现错误：", err), c.Player)
	}
	retStr := string(buffer)
	/**
	插件返回数据以空格区分参数
	第一个为调用方法名
	第二个为方法参数
	第三个如果有则代表玩家名
	*/
	argv := strings.Fields(retStr)
	if len(argv) >= 2 {
		switch argv[0] {
		case "say":
			s.Say(argv[1])
		case "tell":
			if len(argv) >= 3 {
				s.Tell(argv[1], argv[2])
			}
		case "Execute":
			s.Execute(argv[1])
		}
	}
}
