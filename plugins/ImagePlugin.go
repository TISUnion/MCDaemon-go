/**
 * 开启、关闭镜像
 * 前置插件： 备份插件backup
 */

package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/container"
	"MCDaemon-go/lib"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-ini/ini"
)

type ImagePlugin struct{}

func (lp *ImagePlugin) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) == 0 {
		c.Argv = append(c.Argv, "help")
	}
	cor := container.GetInstance()
	switch c.Argv[0] {
	case "show":
		backupfiles, _ := filepath.Glob("back-up/*")
		//标记是否启动
		for k, _ := range backupfiles {
			var split string
			if runtime.GOOS == "windows" {
				split = "\\"
			} else {
				split = "/"
			}
			backupfiles[k] = strings.Split(backupfiles[k], split)[1]
			if cor.IsRuntime(backupfiles[k]) {
				backupfiles[k] += "   已启动" + "   端口：" + cor.Servers[backupfiles[k]].GetPort()
			} else {
				backupfiles[k] += "   未启动"
			}
		}
		text := "备份如下：\\n" + strings.Join(backupfiles, "\\n")
		s.Tell(c.Player, text)
	case "start":
		if len(c.Argv) == 1 {
			s.Tell(c.Player, "缺少启动的镜像名称")
		}
		//镜像已启动则不启动镜像
		if cor.IsRuntime(c.Argv[1]) {
			s.Tell(c.Player, "该镜像已经启动")
		} else {
			path := "back-up/" + c.Argv[1] + "/server.properties"
			if _, err := os.Stat(path); err != nil {
				s.Tell(c.Player, "镜像不存在")
			} else {
				cor = container.GetInstance()
				svr := s.Clone()
				//修改端口
				sercfg, _ := ini.Load(path)
				sercfg.Section("").NewKey("server-port", svr.GetPort())
				sercfg.SaveTo(path)
				//启动
				cor.Add(c.Argv[1], "back-up/"+c.Argv[1], svr)
			}
		}
	case "stop":
		if len(c.Argv) == 1 {
			s.Tell(c.Player, "缺少停止的镜像名称")
		}
		cor = container.GetInstance()
		if cor.IsRuntime(c.Argv[1]) {
			cor.Del(c.Argv[1])
		} else {
			s.Tell(c.Player, "镜像未启动")
		}
	default:
		text := "!!image show 查看镜像\\n!!image start [镜像名称] 开启镜像 \\n!!image stop [镜像名称] 关闭镜像"
		s.Tell(c.Player, text)
	}
}

func (lp *ImagePlugin) Init(s lib.Server) {
}

func (hp *ImagePlugin) Close() {
}
