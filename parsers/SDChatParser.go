package parser

import (
	"MCDaemon-go/command"
	"regexp"
)

type SDChatParser struct {
	startChatPlayer map[string]bool
}

func (sp *SDChatParser) Parsing(word string) (*command.Command, bool) {
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
		//如果已开启聊天模式(私聊模式)
		if command.Group.HasPlayer("SDChat", result["player"]) {
			return &command.Command{
				Player: result["player"],
				Argv:   []string{"say", result["word"]},
				Cmd:    "!!SDChat",
			}, true
		}
		//如果已开启聊天模式(全局聊天模式)
		if command.Group.HasPlayer("SDChat-all", result["player"]) {
			return &command.Command{
				Player: result["player"],
				Argv:   []string{"say-all", result["word"]},
				Cmd:    "!!SDChat",
			}, true
		}
	}
	return nil, false
}
