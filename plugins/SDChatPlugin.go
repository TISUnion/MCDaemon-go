package plugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/config"
	"MCDaemon-go/lib"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

type SDChatPlugin struct{}

func (hp *SDChatPlugin) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) < 1 {
		c.Argv = append(c.Argv, "help")
	}
	switch c.Argv[0] {
	case "all":
		command.Group.AddPlayer("SDChat-all", c.Player)
		s.Tell(c.Player, "开启全局聊天模式成功")
	case "start":
		command.Group.AddPlayer("SDChat", c.Player)
		s.Tell(c.Player, "开启聊天模式成功")
	case "stop":
		command.Group.DelPlayer("SDChat", c.Player)
		command.Group.DelPlayer("SDChat-all", c.Player)
		s.Tell(c.Player, "退出聊天模式成功")
	case "say":
		s.Tell(c.Player, "沙雕："+chat(c.Argv[1], c.Player))
	case "say-all":
		s.Say("沙雕对：" + c.Player + "说" + chat(c.Argv[1], c.Player))
	case "reload":
		_ = config.GetPluginCfg(true)
		s.Tell(c.Player, "已重新读取配置文件")
	default:
		text := "!!SDChat all start 开启全局聊天模式\\n!!SDChat start 开启私聊模式（别的玩家看不见沙雕机器人给你发的信息）\\n!!SDChat stop 关闭聊天模式"
		s.Tell(c.Player, text)
	}
}

func (hp *SDChatPlugin) Init(s lib.Server) {
}

func (hp *SDChatPlugin) Close() {
}

//封装JSON
func LightEncode(elememt interface{}) string {
	//拼接的结果字符串
	var s string
	//若为对象，则拼接字符串
	if LightJson, err := elememt.(map[string]interface{}); !err {
		s = string("{")
		for key, val := range LightJson {
			s += fmt.Sprintf("\"%s\":\"%s\",", key, LightEncode(val))
		}
		s += string("}")
	} else {
		jsonStr, err := json.Marshal(elememt)
		if err != nil {
			log.Fatal("Can't transform jsonString,Because ", err)
		}
		s = string(jsonStr)
	}
	return s
}

//http POST请求
func httpPost(data string) string {
	resp, err := http.Post("http://openapi.tuling123.com/openapi/api/v2",
		"application/x-www-form-urlencoded",
		strings.NewReader(data),
	)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return string(body)
}

func stringToInt(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

//发出请求获取聊天回复
func chat(data string, player string) string {
	_requestMap := map[string]interface{}{
		"perception": map[string]interface{}{
			"inputText": map[string]interface{}{
				"text": data,
			},
		},
		"userInfo": map[string]interface{}{
			"apiKey":     config.GetPluginCfg(false).Section("SDChat").Key("appid").String(),
			"userId":     stringToInt(player),
			"groupId":    10,
			"userIdName": player,
		},
	}
	return gjson.Get(httpPost(LightEncode(_requestMap)), "results.0.values.text").String()
}
