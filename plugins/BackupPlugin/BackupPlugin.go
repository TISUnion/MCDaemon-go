package BackupPlugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/config"
	"MCDaemon-go/lib"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type BackupPlugin struct {
	backupName string
}

func (bp *BackupPlugin) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) == 0 {
		c.Argv = append(c.Argv, "help")
	}
	switch c.Argv[0] {
	case "save":
		if len(c.Argv) < 2 {
			s.Tell(c.Player, "请输入备份存档名称！")
		}
		bp.backupName = c.Argv[1]
		s.Execute("/save-all flush")
	case "saved":
		server_path := config.Cfg.Section("MCDeamon").Key("server_path").String()
		if err := Copy(server_path, "back-up/"+bp.backupName); err != nil {
			fmt.Println(err)
		}
		s.Say("备份完成")
	case "compress":
		if runtime.GOOS == "windows" {
			s.Tell(c.Player, "windows服务器不支持压缩功能")
		} else {
			cmd := exec.Command("tar", "zcvf", fmt.Sprintf("back-up/%s/%s.tar.gz", bp.backupName, bp.backupName), "back-up/"+bp.backupName)
			if err := cmd.Run(); err != nil {
				s.Tell(c.Player, fmt.Sprint("压缩姬出问题了，因为", err))
			} else {
				s.Tell(c.Player, "压缩备份完成")
			}
		}
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
		}
		text := "备份如下：\\n" + strings.Join(backupfiles, "\\n")
		s.Tell(c.Player, text)
	default:
		text := "使用规则：\\n!!backup save [存档名称]\\nlinux下可使用!!backup compress对最近一次save的存档进行压缩\\n!!backup show查看已备份列表"
		s.Tell(c.Player, text)
	}
}

func (bp *BackupPlugin) Init(s lib.Server) {
}

func (bp *BackupPlugin) Close() {
}
