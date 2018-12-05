/**
 *命令语法解析器
 */
package parse

import "MCDaemon-go/lib"

type ParseMachine interface {
	Parsing(string) (*lib.Command, bool)
}

//语法解析器队列
type ParseMachines []ParseMachine
