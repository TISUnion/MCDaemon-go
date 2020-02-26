/*
 * 自动定时备份插件
 * 定时运行：rsync -a --delete minecraft/ back-up/auto
 * 他会对所有人说：嘤嘤嘤
 */

 package plugin

 import 
	"MCDaemon-go/command"
	"MCDaemon-go/lib"
	"fmt"
)
 
 type Yinyin struct {
 }
 
 func (hp *Yinyin) Handle(c *command.Command, s lib.Server) {
	 s.Say(fmt.Sprintf("%s对所有人说：嘤嘤嘤！", c.Player))
 }
 
 func (hp *Yinyin) Init(s lib.Server) {
 }
 
 func (hp *Yinyin) Close() {
 }
 