package parser

import (
	"MCDaemon-go/command"
	"regexp"
)

type WarnParser struct{}

func (p *WarnParser) Parsing(word string) (*command.Command, bool) {
	re := regexp.MustCompile(`\[\d+:\d+:\d+\]\s+\[Server thread/WARN\]: Can't keep up! ` +
		`Is the server overloaded\? Running (\d*)+ms or (\d*)+ ticks behind`)
	match := re.FindStringSubmatch(word)
	if len(match) == 3 {
		_commond := &command.Command{
			Cmd:  "!!Warn",
			Argv: []string{"warn", match[1], match[2]},
		}
		return _commond, true
	}

	return nil, false
}
