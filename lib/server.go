package lib

import "MCDaemon-go/command"

type Server interface {
	Say(string)
	Tell(string, string)
	Execute(string)
	Close()
	Restart()
	RunPlugin(*command.Command)
	RunUniquePlugin(func())
}
