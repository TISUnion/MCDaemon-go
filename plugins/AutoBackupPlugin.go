/*
 * 自动定时备份插件（仅限Linux，在Ubuntu18.04通过测试）
 * 定时运行：rsync -a --delete minecraft/ back-up/auto
 * 他会对所有人说：嘤嘤嘤
 */

package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/lib"
	"fmt"
	"os"
	"time"
)

type AutoBackup struct {
}

func (hp *AutoBackup) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) == 0 {
		c.Argv = append(c.Argv, "help")
	}

	switch c.Argv[0] {
	case "set":
		if len(c.Argv) < 2 {
			s.Tell(c.Player, command.Text{"请输入要设定的小时数！", "red"})
		} else {
			s.Tell(c.Player, command.Text{fmt.Sprintf("收到参数%s", c.Argv[1]), "green"})
		}
		s.Tell(c.Player, command.Text{"本插件还未完工！", "red"})

	case "query":
		lastTime, strerr := GetFileModTime("back-up/auto")
		currentTime := time.Now()
		printTime := "上次存档："
		if len(strerr) > 0 {
			s.Tell(c.Player, command.Text{strerr, "red"})
		} else {
			if currentTime.Year() != lastTime.Year() {
				printTime += fmt.Sprintf("%d年", lastTime.Year())
			}

			printTime += fmt.Sprintf("%d月%d日%d时%d分", lastTime.Month(), lastTime.Day(), lastTime.Hour(), lastTime.Minute())
			s.Tell(c.Player, command.Text{printTime, "yellow"})
		}

		s.Tell(c.Player, command.Text{"本插件还未完工！", "red"})

	default:
		set1 := command.Text{"!!autobk set [小时数]", "white"}
		set2 := command.Text{"设定存档间隔\n", "green"}
		query1 := command.Text{"!!autobk query", "white"}
		query2 := command.Text{"查询上次存档时间\n", "green"}
		s.Tell(c.Player, set1, set2, query1, query2)
	}
}

// GetFileModTime ：获取文件修改时间 返回时间
func GetFileModTime(path string) (t time.Time, strerr string) {
	f, err := os.Open(path)
	if err != nil {
		return time.Now(), "open file error"
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return time.Now(), "stat fileinfo error"
	}

	return fi.ModTime(), ""
}

func (hp *AutoBackup) Init(s lib.Server) {
	fmt.Println("ooooooh!")
}

func (hp *AutoBackup) Close() {
}
