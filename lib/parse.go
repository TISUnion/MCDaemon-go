package lib

import "reloadMCDaemon-go/command"

type ParseMachine interface {
	Parsing(string) (*command.Command, bool)
}
