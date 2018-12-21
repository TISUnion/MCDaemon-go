package lib

import (
	"MCDaemon-go/command"
)

type Server interface {
	Say(...interface{})
	Tell(string, ...interface{})
	Execute(string)
	Close()
	Restart()
	Getinfo() string
	Clone(string, []string) Server
	RunPlugin(*command.Command)
	RunUniquePlugin(func())
	WriteLog(level string, msg string)
	GetPluginList() map[string]Plugin
	GetDisablePluginList() map[string]Plugin
	GetParserList() []Parser
}
