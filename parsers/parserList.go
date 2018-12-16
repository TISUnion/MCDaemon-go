package parser

import "MCDaemon-go/lib"

var ParseList []lib.Parser

//语法解析器列表
func init() {
	ParseList = []lib.Parser{
		defaultParser{},
	}
}
