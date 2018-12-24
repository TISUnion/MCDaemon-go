/*
 * 插件编写例子
 * 当玩家输入命令：!!yinyinmaster
 * 他会对所有人说：嘤嘤嘤
 */

package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/lib"
	"fmt"
)

type Yinyin struct {
}

func (hp *Yinyin) Handle(c *command.Command, s lib.Server) {
	s.Say(fmt.Sprintf("%s对所有人说：嘤嘤嘤！", c.Player))
}

func (hp *Yinyin) Init(s lib.Server) {
}
