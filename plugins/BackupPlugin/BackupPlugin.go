package BackupPlugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/config"
	"MCDaemon-go/lib"
	"fmt"
	"os/exec"
	"runtime"
)

type BackupPlugin struct {
	backupName string
}

func (bp *BackupPlugin) Handle(c *command.Command, s lib.Server) {
	server_path := config.Cfg.Section("MCDeamon").Key("server_path").String()
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
		if err := Copy(server_path, "back-up/"+bp.backupName); err != nil {
			fmt.Println(err)
		}
		s.Say("备份完成")
	case "compress":
		if runtime.GOOS == "windows" {
			s.Tell(c.Player, "windows服务器不支持压缩功能")
		} else {
			cmd := exec.Command("tar", "zcvf", "back-up/"+bp.backupName+".tar.gz", "back-up/"+bp.backupName)
			if err := cmd.Run(); err != nil {
				s.Tell(c.Player, fmt.Sprint("压缩姬出问题了，因为", err))
			} else {
				s.Tell(c.Player, "备份完成")
			}
		}
	default:
		text := "使用规则：\\n!!backup save [存档名称]\\nlinux下可使用!!backup compress对最近一次save的存档进行压缩"
		s.Tell(c.Player, text)
	}
}
