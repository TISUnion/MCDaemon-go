package parser

import (
	"MCDaemon-go/command"
	"regexp"
)

type ChatParser struct {
	startChatPlayer map[string]bool
}

func (sp *ChatParser) Parsing(word string) (*command.Command, bool) {
	if sp.startChatPlayer == nil {
		sp.startChatPlayer = make(map[string]bool)
	}
	//匹配玩家说的话
	re := regexp.MustCompile(`\[\d+:\d+:\d+\]\s+\[Server thread/INFO\]:\s+<(?P<player>.+)>\s+(?P<word>.+)\s*`)
	match := re.FindStringSubmatch(word)
	groupNames := re.SubexpNames()
	result := make(map[string]string)

	//是否是玩家说的话
	if len(match) != 0 {
		// 转换为map
		for i, name := range groupNames {
			result[name] = match[i]
		}
		//如果已开启聊天模式(全局聊天模式)
		if command.Group.HasPlayer("ServersChat", result["player"]) {
			return &command.Command{
				Player: result["player"],
				Argv:   []string{"chat_xxx_say", result["word"]},
				Cmd:    "!!Chat",
			}, true
		}
	}
	return nil, false
}
