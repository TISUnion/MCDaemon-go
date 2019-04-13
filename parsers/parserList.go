package parser

import "MCDaemon-go/lib"

//语法解析器列表
func CreateParserList() []lib.Parser {
	return []lib.Parser{
		&defaultParser{},
		&BackupParser{},
		&SDChatParser{},
		&TpsParser{},
	}
}
