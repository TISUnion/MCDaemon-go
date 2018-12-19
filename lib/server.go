package lib

import (
	"MCDaemon-go/command"
)

type Server interface {
	Say(string)
	Tell(string, string)
	Execute(string)
	Close()
	Restart()
	Clone(string, []string) Server
	RunPlugin(*command.Command)
	RunUniquePlugin(func())
	WriteLog(level string, msg string)
	GetPluginList() map[string]Plugin
	GetDisablePluginList() map[string]Plugin
	GetParserList() []Parser
}
