/*
 * 自动增量备份插件（仅限Linux，在Ubuntu18.04通过测试）
 * 定时运行：rsync -a --delete minecraft/ back-up/auto
 */

package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/config"
	"MCDaemon-go/lib"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

type AutoBackup struct {
	interval int
	workdir  string
}

func (ab *AutoBackup) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) == 0 {
		c.Argv = append(c.Argv, "help")
	}

	switch c.Argv[0] {
	case "set":
		if len(c.Argv) < 2 {
			s.Tell(c.Player, command.Text{"请输入要设定的小时数！", "red"})
		} else {
			interval_new, err := strconv.Atoi(c.Argv[1])
			if err != nil || interval_new < 0 {
				s.Tell(c.Player, command.Text{"请输入自然数！", "red"})
			} else if interval_new == 0 {
				s.Say(command.Text{"自动存档关闭！", "yellow"})
				ab.interval = interval_new
			} else {
				s.Say(command.Text{fmt.Sprintf("自动存档间隔设为%d小时", interval_new), "yellow"})
				ab.interval = interval_new
			}
		}

	case "query":
		lastTime, strerr := GetFileChangeTime("back-up/auto")
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
		if ab.interval > 0 && len(ab.workdir) > 0 {
			s.Tell(c.Player, command.Text{fmt.Sprintf("间隔为%d小时", ab.interval), "yellow"})
		} else {
			s.Tell(c.Player, command.Text{"自动备份已关闭", "yellow"})
		}

	case "save":
		if ab.interval > 0 && len(ab.workdir) > 0 {
			lastTime, strerr := GetFileChangeTime("back-up/auto")
			currentTime := time.Now()
			if len(strerr) > 0 {
				s.Tell(c.Player, command.Text{strerr, "red"})
			} else if currentTime.Unix()-lastTime.Unix() > int64(ab.interval*3600) {
				s.Say(command.Text{"开始自动备份...", "yellow"})
				// 将上次备份转存
				cmdlast := exec.Command("rsync", "-a", "--delete", "back-up/auto", "back-up/auto-last")
				errlast := cmdlast.Run()
				if errlast != nil {
					s.Say(command.Text{fmt.Sprintf("备份错误：%s", errlast), "red"})
					fmt.Println(fmt.Sprintf("备份错误：%s", errlast))
				} else {
					// 备份当前存档
					cmd := exec.Command("rsync", "-a", "--delete", ab.workdir+"/", "back-up/auto")
					err := cmd.Run()
					if err != nil {
						s.Say(command.Text{fmt.Sprintf("备份错误：%s", err), "red"})
						fmt.Println(fmt.Sprintf("备份错误：%s", err))
					} else {
						s.Say(command.Text{"备份完毕", "green"})
						fmt.Println("备份完毕")
					}
				}

			}
		}

	default:
		set1 := command.Text{"!!autobk set [小时数]", "white"}
		set2 := command.Text{"设定存档间隔，0为关闭\n", "green"}
		query1 := command.Text{"!!autobk query", "white"}
		query2 := command.Text{"查询上次存档时间和间隔时间\n", "green"}
		s.Tell(c.Player, set1, set2, query1, query2)
	}
}

// GetFileChangeTime ：获取文件修改时间 返回时间
func GetFileChangeTime(path string) (t time.Time, strerr string) {
	f, err := os.Open(path)
	if err != nil {
		return time.Now(), "open file error"
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return time.Now(), "stat fileinfo error"
	}

	filestat := fi.Sys().(*syscall.Stat_t)
	return timespecToTime(filestat.Atim), ""
}

// timespecToTime ：将获取到的元信息时间转为Time
func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func (ab *AutoBackup) Init(s lib.Server) {
	ab.interval = 1
	ab.workdir = config.GetPluginCfg(false).Section("AutoBackup").Key("workdir").String()
}

func (ab *AutoBackup) Close() {
}
