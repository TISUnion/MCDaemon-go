# MCDaemon-go

## 用golang实现的Minecraft进程管理程序

在windows10以及centos7中运行成功

-----
## 开始使用
- 下载最新的release(beta版不提供)
- 最新版的release默认包含了[插件收录库](https://github.com/TISUnion/MCDaemonPlugins-go)的所有插件
### 快速开始
- 修改配置文件MCD_conig.ini
  - 解压最新版MCDaemon,进入并创建一个minecraft文件夹
  - 将下载MC服务端放入创建的minecraft文件夹内，重命名为server.jar
- 运行start(linux/unix)或者start.exe（windows）
- 在游戏中输入!!server show命令来查看可使用的命令
-----
## 配置文件
### MCD_conig.ini
- run_environment：运行方式默认为develop开发者模式
- server_name&server_path：决定了服务器启动后生成文件的位置，服务器文件的logs, world等文件都会生成在填写的server_path路径中，默认为`minecraft`。server_name则是要运行服务端文件名，默认为`server.jar`
- Xms： jvm运行的最小堆内存
- Xmx：jvm最大堆内存
- gui：是否使用图形界面，默认为true。`注意`,只有带有图形界面的操作系统才会起作用。
- maxRunPlugins：插件最大并发数，默认为`10`，所有插件都是异步运行的，输出是同步输出到服务端。
- plugins: 该域下填写的都是热插件参数，键为调用插件的命令，值为插件对应的二进制文件名称，`注意`, 二进制文件都在hotPlugins文件夹下
-----
## 插件编写

- ### 热插件
   - 将热插件的执行文件放入hotPlugins文件夹中
   - 在MCD_conig.ini的[plugins]域中注册热插件
   - 通过如下代码获取插件的命令的参数：
     ```golang
        args := os.Args
        args = args[1:]
     ```
   - 执行文件的输出：第一个参数固定为回调函数名称，参数之间用空格区别，提供三个函数，三个函数分别为：
   - say: 对全服玩家说话，第二个参数为说的字符串
   - tell: 对某一玩家说话，第二个参数为玩家名称， 第二个参数为说的字符串
   - Execute: 执行某一命令，第二个参数为命令字符串

   输出例子：
   ```text
   say hello everyone  //对所有人说hello everyone
   tell Alice hello
   Execute debug start
   ```
- ### 冷插件
   所有冷插件都在plugins文件夹中，包含一个栗子和一个基础插件。
   
   所有插件需要实现lib包里的Plugin接口，Handle方法为插件的调用方法，当输入命令触发插件时，就会运行该方法。
   例子：
   ```golang
   func (hp *Yinyin) Handle(c *command.Command, s lib.Server) {
	   s.Say(fmt.Sprintf("%s对所有人说：嘤嘤嘤！", c.Player)) //调用lib.Server接口的say方法对全服人说话
   }
   ```
   传入的参数为命令对象(command.Command)和被调用的服务器的接口(lib.Server)，命令对象中有解释器传过来的玩家名称，调用命令以及调用命令参数等属性。服务器接口则提供插件开发需要所有接口。服务器接口的方法会在下面详细讲述。
   
   `注意` 开发完成的插件需要到插件列表中注册，在plugins下的
   pluginList.go中添加如下代码：
   ```golang
   ....
   //注册冷插件
   PluginsList.RegisterPlugin("命令", &PluginName{})                //注册XXX插件
   ....
   ```
   这里的命令就是游戏中输入的命令，默认是!或者!!作为前缀，如果想修改可以直接修改语法解释器或者自定义一个新的语法解释器

   如果想要获取插件配置文件，可以使用config.GetPluginCfg()方法，该方法返回一个*ini.File类型，具体用法可以参考[go-ini](https://github.com/go-ini/ini)这个第三方库,插件配置文件是Plugin_conf.ini，配置参数写在一个域内。
- ### lib.Server接口方法
   - Say 对全服发送消息
   - Tell 对某一个玩家发送消息，第二个参数玩家昵称，第三个为发送的消息字符串
   - Execute 执行一个原版MC命令
   - Close 关闭当前服务器实例对应服务器
   - Restart 重启当前服务器实例对应服务器
   - Start 开启当前服务器实例里配置的服务器
   - ReloadConf 重新读取配置文件，加入热插件后需要重新读取配置文件才能生效
   - GetPluginList 获取可使用插件命令列表
   - GetDisablePluginList 获取被禁用插件列表
   - GetParserList 获取解释器插件列表
   - RunPlugin 执行插件传入一个command.Command对象，即可在一个插件中调用其他插件
   - RunUniquePlugin 传入一个函数，等待所有正在运行的插件运行完毕后堵塞住插件运行池，然后运行传入的函数。
   - WriteLog 写入日志，第一个参数为消息等级，第二个为写入日志的消息，日志可以在logs文件中查看，日志登级分为：debug、info、warn、error和fatal。
   - Clone 克隆出一个server.Server实例用于承载镜像
   - GetPort 开启服务器的端口，配置文件中为默认端口，镜像则会随机一个可用的端口。
   - Getinfo 获取镜像名称和端口
   - CloseInContainer 在容器中关闭镜像，与Close不同的是Close用单服务器，无法注销容器中对应的服务器数据

   以上方法除了Close都适用于单服务器和多服务器运行，后面会提到多服务器镜像的使用

-----   

## 高级组件
- ### 自定义语法解释器
   所有语法解释器都在parsers中，所有解释器需要实现lib中的Parser接口。
   
   Parsin方法会传入一个字符串，是从服务器接受的一条信息，例如玩家发言的结构为：
   
   [时间] [Server thread/INFO]: 玩家名 说的话
   
   以及执行命令完服务器返回的特殊结构的数据(一般都为json格式)，所以我们需要对该字符串做强正则，因为程序会对该条信息运行所有的语法解释器，需保证该解释器不会和其他解释器捕获同一条消息，当然如果想要两个解释器捕获同一条信息也是可以的。
   
   和插件类似，写完的解释器需要在parsers文件夹下的parserList.go中加入实现的解释器对象。

   例子：
   ```golang
   ....
   return []lib.Parser{
		....
		XXXParser{},
		....
	}
   ```
- ### 多服务器镜像

   镜像主要基于备份插件，所以把备份插件和镜像插件合并进了主分支。如果要使用多镜像，需要了解容器，容器代码在container包下的Container.go中，下面介绍如何使用容器：
   - container.GetInstance() 获取全局容器实例，这是一个单例，容器原则上整个程序运行时只允许有一个。开发者也可以自己在插件中创建容器实例，但是并不推荐。
   - Add 容器实例的方法，添加一个服务器镜像。第一个参数为镜像名称，用于区分镜像，第二各参数为镜像所在的目录，第三个参数为lib.server接口，用改Clone一个Server实例来承载服务器。
   - Del 关闭一个服务器镜像，传入镜像名称即可
   - GetRuntimeServer 获取所有正在运行镜像的lib.Server接口，可用来说像所有镜像服务器发送消息甚至运行插件。
   - IsRuntime 传入一个镜像名称，查询改镜像是否启动

- ### 玩家组
   在对玩家进行分类分组时，由于每个插件在同一个服务器中的实例就一个所以，在单服务器中可以自建玩家组区分玩家，但是如果想对全部镜像进行玩家分组，则需要使用提供玩家组组件。

   玩家组代码在command包的playerGroup.go里，使用方法如下：

   - command.Group 全局玩家组实例，在该实例下有一系列创建玩家组，在组中添加删除玩家的方法
   - AddPlayer 玩家组实例方法, 将玩家加入到指定用户组中, 第一个参数为用户组名称,第二个参数为玩家名称
   - DelPlayer 将玩家在指定用户组中删除, 第一个参数为用户组名称,第二个参数为玩家名称
   - HasPlayer 查询玩家在指定用户组中是否存在, 第一个参数为用户组名称,第二个参数为玩家名称
   - GetPlayer 获取整个玩家组数组
- ### 色彩缤纷的消息回复
   在使用lib.Server接口的Say以及Tell方法时, 不仅仅可以传入字符串, 还可能传入command.Text类型的数据, command.Text对象只有2个参数, 一个是消息字符串, 一个消息的颜色, 可以传入多个Text来生成一句各种颜色的消息。
   
   例子：
   ```golang
   s.Say(command.Text{"hello", "green"}, command.Text{"everyone", "red"})

   s.Tell(playerName, command.Text{"hello", "green"}, command.Text{playerName, "red"})

   //也可以传入数组
   s.Say([]command.Text{command.Text{"hello", "green"}, command.Text{"everyone", "red"}})

   s.Tell(playerName, []command.Text{command.Text{"hello", "green"}, command.Text{playerName, "red"}})
   ```