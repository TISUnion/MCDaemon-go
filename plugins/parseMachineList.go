package plugin

import (
	"MCDaemon-go/parse"
)

var ParseMachineList parse.ParseMachines

//加载所有语法分析器
func init() {
	ParseMachineList = parse.ParseMachines{
		parse.ParseCommand{},
	}
}
