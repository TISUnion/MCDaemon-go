package parser

import (
	"MCDaemon-go/command"
	"regexp"
)

type HereParser struct{}

func (p *HereParser) Parsing(word string) (*command.Command, bool) {
	//re := regexp.MustCompile(`\[\d+:\d+:\d+\]\s+\[Server thread/INFO\]:\s(\S*)\shas the following entity data`)
	//player := re.FindStringSubmatch(word)
	//if len(player) == 2 {
		re := regexp.MustCompile(`Dimension:\s(\d)`)
		dimension := re.FindStringSubmatch(word)
		if len(dimension) == 2 {
			re = regexp.MustCompile(`Pos:\s\[([-0-9]*).[-0-9]*d,\s([-0-9]*).[-0-9]*d,\s([-0-9]*).[-0-9]*d\]`)
			pos := re.FindStringSubmatch(word)
			if len(pos) == 4 {
				_commond := &command.Command{
					Cmd:  "!!here",
					Argv: []string{"res", dimension[1], pos[1], pos[2], pos[3]},
				}
				return _commond, true
			}
		}
	//}
	return nil, false
}
