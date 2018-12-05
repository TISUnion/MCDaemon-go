package config

import (
	"MCDaemon-go/lib"
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

type Config struct {
	Conf    map[string]string
	Plugins map[string]lib.HotPlugin
}

//配置变量
var (
	cfg      *ini.File
	err      error
	is_build bool
	conf     *Config
)

func (c *Config) GetStartConfig() map[string]string {
	if c.Conf == nil {
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
		c.Conf = map[string]string{
			"serverName": serverName,
			"serverPath": serverPath,
			"Xmx":        xmx,
			"Xms":        xms,
			"gui":        gui,
		}
	}
	return c.Conf
}

func (c *Config) GetPlugins() map[string]lib.HotPlugin {
	if c.Plugins == nil {
		cfg, _ = ini.Load("MCD_conig.ini")
		plugins := make(map[string]lib.HotPlugin)
		keys := cfg.Section("plugins").KeyStrings()
		for _, val := range keys {
			plugins[val] = lib.HotPlugin(cfg.Section("plugins").Key(val).String())
		}
		c.Plugins = plugins
	}
	return c.Plugins
}

func GetInstance() *Config {
	if !is_build {
		conf = &Config{}
	}
	return conf
}

func init() {
	is_build = false
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
