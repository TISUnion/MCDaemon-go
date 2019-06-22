package lib

import (
	"MCDaemon-go/config"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	devlog     *logrus.Logger
	commonlog  *logrus.Logger
	isdevelop  bool
	printToCmd bool
)

func init() {
	devlog = logrus.New()
	commonlog = logrus.New()
	devlog.SetLevel(logrus.DebugLevel)
	commonlog.SetLevel(logrus.DebugLevel)
	if config.Cfg.Section("").Key("run_environment").String() == "develop" {
		isdevelop = true
	} else {
		isdevelop = false
	}
	//开发时调试使用，是否显示在命令行中
	printToCmd = false
}

//写入开发日志
func WriteDevelopLog(level string, msg string) {
	//如果不是开发者模式，则不写入日志
	if !isdevelop {
		return
	}
	if printToCmd {
		fmt.Println(msg)
	}
	logFile, err := os.OpenFile("logs/develop.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	defer logFile.Close()
	if err != nil {
		fmt.Println("日志写入系统发生错误！ 因为", err)
	}
	devlog.SetOutput(logFile)
	switch level {
	case "debug":
		devlog.Debug(msg)
	case "info":
		devlog.Info(msg)
	case "warn":
		devlog.Warn(msg)
	case "error":
		devlog.Error(msg)
	case "fatal":
		devlog.Fatal(msg)
	}
}

//写入运行日志
func WriteRuntimeLog(level string, msg string, serverName string) {
	logFile, err := os.OpenFile(fmt.Sprintf("logs/%s_%s.log", serverName, time.Now().Format("2006-01-02")), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	defer logFile.Close()
	if err != nil {
		fmt.Println("日志写入系统发生错误！ 因为", err)
	}
	commonlog.SetOutput(logFile)
	switch level {
	case "debug":
		commonlog.Debug(msg)
	case "info":
		commonlog.Info(msg)
	case "warn":
		commonlog.Warn(msg)
	case "error":
		commonlog.Error(msg)
	case "fatal":
		commonlog.Fatal(msg)
	}
}
