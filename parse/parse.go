/**
 *命令语法解析器
 */
package parse

import "MCDaemon-go/command"

type ParseMachine interface {
	Parsing(string) (*command.Command, bool)
}

//语法解析器队列
type ParseMachines []ParseMachine
