package parser

import (
	"MCDaemon-go/command"
	"regexp"
)

type LoginoutParser struct{}

func (p *LoginoutParser) Parsing(word string) (*command.Command, bool) {
	re := regexp.MustCompile(`\[\d+:\d+:\d+\]\s+\[Server thread/INFO\]:\s(\S*)\sleft the game`)
	match := re.FindStringSubmatch(word)
	if len(match) == 2 {
		_commond := &command.Command{
			Cmd:  "!!autobk",
			Argv: []string{"save", match[1]},
		}
		return _commond, true
	}

	// re = regexp.MustCompile(`\[\d+:\d+:\d+\]\s+\[Server thread/INFO\]:\s(\S*)\sjoined the game`)
	// match = re.FindStringSubmatch(word)
	// if len(match) == 2 {
	// 	_commond := &command.Command{
	// 		Cmd:  "!!autobk",
	// 		Argv: []string{"save"},
	// 	}
	// 	return _commond, true
	// }

	return nil, false
}
