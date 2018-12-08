package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/lib"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

//热加载插件类型
type HotPlugin string

func (hp HotPlugin) Handle(c *command.Command, s lib.Server) {
	commandName := "./hotPlugins/" + c.PluginName
	pluginProcess := exec.Command(commandName, c.Argv...)
	stdout, _ := pluginProcess.StdoutPipe()
	defer stdout.Close()
	if err := pluginProcess.Start(); err != nil {
		s.Tell(fmt.Sprintf("%s插件出错！ 因为%v", c.PluginName, err), c.Player)
	}
	buffer, err := ioutil.ReadAll(stdout)
	if err != nil {
		if err != io.EOF {
			msg := fmt.Sprintf("%s插件出错！ 因为%v", c.PluginName, err)
			s.Tell(msg, c.Player)
		}
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
