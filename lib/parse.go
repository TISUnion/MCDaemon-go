package lib

import "MCDaemon-go/command"

type Parser interface {
	Parsing(string) (*command.Command, bool)
}
