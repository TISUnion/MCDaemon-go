package parser

import (
	"MCDaemon-go/command"
	"regexp"
)

type BackupParser struct{}

func (p BackupParser) Parsing(word string) (*command.Command, bool) {
	re := regexp.MustCompile(`\[\d+:\d+:\d+\]\s+\[Server thread/INFO\]:\s+Saved the game`)
	match := re.FindStringSubmatch(word)
	//匹配tps信息时
	if len(match) != 0 {
		_commond := &command.Command{
			Cmd:  "!!backup",
			Argv: []string{"saved"},
		}
		return _commond, true
	}
	return nil, false
}
