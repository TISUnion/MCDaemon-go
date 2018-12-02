package config

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

//配置变量
var (
	serverName string
	serverPath string
	Xms        string
	Xmx        string
	gui        string
)

func GetConfig() map[string]string {
	//加载配置文件
	cfg, err := ini.Load("MCD_conig.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	//读取配置
	Section := cfg.Section("MCDeamon")
	serverName = Section.Key("server_name").String()
	serverPath = Section.Key("server_path").String()
	//设置默认值
	Xms = Section.Key("Xms").Validate(func(in string) string {
		if len(in) == 0 {
			return "-Xms1024M"
		}
		return fmt.Sprint("-Xms", in)
	})
	Xmx = Section.Key("Xmx").Validate(func(in string) string {
		if len(in) == 0 {
			return "-Xmx1024M"
		}
		return fmt.Sprint("-Xmx", in)
	})
	gui = Section.Key("gui").Validate(func(in string) string {
		if len(in) == 0 {
			return "false"
		}
		return in
	})
	return map[string]string{
		"serverName": serverName,
		"serverPath": serverPath,
		"Xmx":        Xmx,
		"Xms":        Xms,
		"gui":        gui,
	}
}
