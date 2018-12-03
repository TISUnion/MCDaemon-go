package config

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

//配置变量
var (
	cfg *ini.File
	err error
)

func GetConfig() map[string]string {
	//加载配置文件
	cfg, err = ini.Load("MCD_conig.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	//读取配置
	Section := cfg.Section("MCDeamon")
	serverName := Section.Key("server_name").String()
	serverPath := Section.Key("server_path").String()
	//设置默认值
	xms := Section.Key("Xms").Validate(func(in string) string {
		if len(in) == 0 {
			return "-Xms1024M"
		}
		return fmt.Sprint("-Xms", in)
	})
	xmx := Section.Key("Xmx").Validate(func(in string) string {
		if len(in) == 0 {
			return "-Xmx1024M"
		}
		return fmt.Sprint("-Xmx", in)
	})
	gui := Section.Key("gui").Validate(func(in string) string {
		if len(in) == 0 {
			return "false"
		}
		return in
	})
	return map[string]string{
		"serverName": serverName,
		"serverPath": serverPath,
		"Xmx":        xmx,
		"Xms":        xms,
		"gui":        gui,
	}
}

func SetEula() {
	cfg, err = ini.Load("eula.txt")
	//不存在eula.txt
	if err != nil {
		cfg = ini.Empty()
		cfg.Section("").NewKey("eula", "true")
		_ = cfg.SaveTo("eula.txt")
	}
	//如果为false
	if cfg.Section("").Key("eula").String() == "false" {
		cfg.Section("").NewKey("eula", "true")
		_ = cfg.SaveTo("eula.txt")
	}
}

func GetPlugins() (map[string]string, map[string]string) {
	cfg, err = ini.Load("MCD_conig.ini")
	userPlugins := make(map[string]string)
	serverPlugins := make(map[string]string)

	keys := cfg.Section("user_plugins").KeyStrings()
	for _, val := range keys {
		userPlugins[val] = cfg.Section("user_plugins").Key(val).String()
	}

	keys = cfg.Section("server_plugins").KeyStrings()
	for _, val := range keys {
		serverPlugins[val] = cfg.Section("server_plugins").Key(val).String()
	}
	return userPlugins, serverPlugins
}
