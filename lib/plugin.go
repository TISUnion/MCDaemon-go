package lib

import "reloadMCDaemon-go/command"

type plugin interface {
	Handle(*command.Command, *Server)
}
