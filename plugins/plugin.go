package plugin

import (
	"MCDaemon-go/lib"
	"MCDaemon-go/server"
)

type Plugin interface {
	Handle(*lib.Command, *server.Server)
}
