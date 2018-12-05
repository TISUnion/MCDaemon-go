package lib

import (
	"MCDaemon-go/server"
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type HotPlugin string

func (hp HotPlugin) Handle(c *Command, s *server.Server) {
	commandName := "./plugin/" + c.PluginName
	pluginProcess := exec.Command(commandName, c.Argv...)
	stdout, _ := pluginProcess.StdoutPipe()
	outRead := bufio.NewReaderSize(stdout, 10000)
	var buffer []byte = make([]byte, 10000)
	if err := pluginProcess.Run(); err != nil {
		s.Tell(c.PluginName+"插件出错！", c.Player)
	}
	n, err := outRead.Read(buffer)
	if err != nil {
		if err != io.EOF {
			msg := fmt.Sprint("%s插件出错！ 因为%v", c.PluginName, err)
			s.Tell(msg, c.Player)
		}
	}
	retStr := string(buffer[:n])
	/**
	以空格区分参数
	第一个为调用方法名
	第二个为方法参数
	第三个如果有则代表玩家名
	*/
	argv := strings.Fields(retStr)
	switch argv[0] {
	case "say":
		s.Say(argv[1])
	case "tell":
		s.Tell(argv[1], argv[2])
	}
}
