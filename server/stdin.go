package server

import (
	"fmt"
	"io"
)

func (svr Server) Say(text string) {
	_command := fmt.Sprintf("/tellraw @a {\"text\":\"%s\"}", text)
	svr.Execute(_command)
}

func (svr Server) Tell(text string, player string) {
	_command := fmt.Sprintf("/tellraw %s {\"text\":\"%s\"}", player, text)
	svr.Execute(_command)
}

func (svr Server) Execute(_command string) {
	//输入的命令要换行！否则无法执行
	_command = _command + "\n"
	_, err := io.WriteString(svr.stdin, _command)
	if err != nil {
		fmt.Println("there is a error!", err)
	}
}
