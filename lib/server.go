package lib

import (
	"MCDaemon-go/command"
)

type Server interface {
	Say(...interface{})
	Tell(string, ...interface{})
	Execute(string)
	Close()
	CloseInContainer()
	Restart()
	Start(string, []string, string)
	Getinfo() string
	Clone() Server
	GetPort() string
	RunPlugin(*command.Command)
	RunUniquePlugin(func())
	WriteLog(level string, msg string)
	GetPluginList() map[string]Plugin
	GetDisablePluginList() map[string]Plugin
	GetParserList() []Parser
}
