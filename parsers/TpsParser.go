package parser

import (
	"MCDaemon-go/command"
	"regexp"
)

type TpsParser struct{}

func (p TpsParser) Parsing(word string) (*command.Command, bool) {
	re := regexp.MustCompile(`\[\d+:\d+:\d+\]\s+\[Server thread/INFO\]:\s+Stopped debug profiling after (?P<TpsInfo>.+)`)
	match := re.FindStringSubmatch(word)
	groupNames := re.SubexpNames()

	result := make(map[string]string)
	//匹配tps信息时
	if len(match) != 0 {
		// 转换为map
		for i, name := range groupNames {
			result[name] = match[i]
		}
		_commond := &command.Command{
			Cmd:  "!!tps",
			Argv: []string{"res", result["TpsInfo"]},
		}
		return _commond, true
	}
	return nil, false
}
