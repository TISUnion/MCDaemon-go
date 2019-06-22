package parser

import "MCDaemon-go/lib"

//语法解析器列表
func CreateParserList() []lib.Parser {
	return []lib.Parser{
		&SDChatParser{},
		&ChatParser{},
		&defaultParser{},
		&BackupParser{},
		&TpsParser{},
	}
}
