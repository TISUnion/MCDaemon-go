package lib

import "MCDaemon-go/command"

type Plugin interface {
	Handle(*command.Command, Server)
	Init(Server)
	Close()
}
