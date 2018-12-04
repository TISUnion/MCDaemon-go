package server

import (
	"MCDaemon-go/parse"
	plugin "MCDaemon-go/plugins"
)

var ParseMachineList parse.ParseMachines

const ParseMachineSize = 10

func init() {
	ParseMachineList = parse.ParseMachines{
		parse.ParseCommand{},
	}
	ParseMachineList = append(ParseMachineList, plugin.ParseMachines...)
}
