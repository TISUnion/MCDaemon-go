package server

import (
	"MCDaemon-go/command"
	"fmt"
	"io"
)

func (svr *Server) Say(argv ...interface{}) {
	svr.Tell("@a", argv...)
}

func (svr *Server) Tell(player string, argv ...interface{}) {
	var stringText string
	var TextArray []command.Text
	var _command string
	for _, v := range argv {
		switch t := v.(type) {
		case string:
			stringText = t
		case []command.Text:
			TextArray = append(TextArray, t...)
		case command.Text:
			TextArray = append(TextArray, t)
		default:
			fmt.Println("不支持的消息类型")
		}
	}
	if stringText != "" {
		_command = fmt.Sprintf("/tellraw %s {\"text\":\"%s\"}", player, stringText)
	} else if len(TextArray) != 0 {
		_command, _ = command.JsonEncode(TextArray)
		_command = fmt.Sprintf("/tellraw %s %s", player, _command)
	}
	svr.Execute(_command)
}

func (svr *Server) Execute(_command string) {
	//输入的命令要换行！否则无法执行
	_command = _command + "\n"
	//同步写入
	svr.lock.Lock()
	defer svr.lock.Unlock()
	_, err := io.WriteString(svr.stdin, _command)
	if err != nil {
		// fmt.Println("there is a error!", err)
	}
}
