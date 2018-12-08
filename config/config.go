package config

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

//配置变量
var (
	cfg     *ini.File
	err     error
	plugins map[string]string
)

//获取服务器启动配置
func GetStartConfig() []string {
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
	result := []string{
		xmx,
		xms,
		"-jar",
		fmt.Sprintf("%s/%s", serverPath, serverName),
		"-Dfile.encoding=UTF-8",
	}
	if gui != "true" {
		result = append(result, "nogui")
	}
	return result
}

//获取插件配置
func GetPlugins(is_rebuild bool) map[string]string {
	if is_rebuild || plugins == nil {
		cfg, _ = ini.Load("MCD_conig.ini")
		plugins = make(map[string]string)
		keys := cfg.Section("plugins").KeyStrings()
		for _, val := range keys {
			plugins[val] = cfg.Section("plugins").Key(val).String()
		}
	}
	return plugins
}

//根据命令获取插件
func GetPluginName(cmd string) string {
	pluins := GetPlugins(false)
	return pluins[cmd]
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
