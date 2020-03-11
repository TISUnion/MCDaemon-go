/*
 * 转发服务器警告消息，并在接收到fatal时自动重启服务器进程
 * author: Sciroccogti
 */

package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/lib"
	"fmt"
	"strconv"
	"time"
)

type WarnPlugin struct{}

func (wp *WarnPlugin) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) != 3 {
		c.Argv[0] = "help"
	}

	switch c.Argv[0] {
	case "warn":
		ticks, _ := strconv.Atoi(c.Argv[2])
		if ticks >= 40 && ticks < 500 {
			s.Say(command.Text{fmt.Sprintf("嗯？服务姬有点忙不过来了，延迟%dticks~", ticks), "grey"})
			s.WriteLog("info", fmt.Sprintf("服务器延迟%dticks", ticks))
		} else if ticks >= 500 && ticks < 1000 {
			s.Say(command.Text{fmt.Sprintf("哎呀呀，让服务姬歇一会吧，延迟%dticks！", ticks), "yellow"})
			s.WriteLog("warn", fmt.Sprintf("服务器延迟%dticks", ticks))
		} else if ticks >= 1000 && ticks < 1500 {
			s.Say(command.Text{fmt.Sprintf("呜呜呜，服务姬受不了了！延迟%dticks！", ticks), "red"})
			s.WriteLog("warn", fmt.Sprintf("服务器延迟%dticks", ticks))
		} else if ticks >= 1500 {
			s.Say(command.Text{"服务器负载过高！请立即停止活动并尽快下线！一分钟后自动重启", "red"})
			s.WriteLog("fatal", fmt.Sprintf("服务器延迟%dticks，即将自动重启", ticks))
			time.Sleep(time.Second * 60)
			s.Restart()
		}
	default:
	}
}

func (wp *WarnPlugin) Init(s lib.Server) {
}

func (wp *WarnPlugin) Close() {
}
