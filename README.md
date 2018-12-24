# MCDaemon-go

## 用golang实现的Minecraft进程管理程序

在windows10以及centos7中运行成功

-----
## 开始使用
- 下载最新的release(beta版不提供)
- 最新版的release默认包含了[插件收录库](https://github.com/TISUnion/MCDaemonPlugins-go)的所有插件
### 快速开始
- 修改配置文件MCD_conig.ini
  1. 解压最新版MCDaemon,进入并创建一个minecraft文件夹
  2. 将下载MC服务端放入创建的minecraft文件夹内，重命名为server.jar
- 根据需求更改配置文件参数
- 运行start(linux/unix)或者start.exe（windows）
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
1. 将热插件的执行文件放入hotPlugins文件夹中
2. 在MCD_conig.ini的[plugins]域中注册热插件
3. 通过如下代码获取插件的命令的参数：
     ```golang
        args := os.Args
        args = args[1:]
     ```
   输出：第一个参数固定为回调函数名称，参数之间用空格区别，提供三个函数，三个函数分别为：
   - say: 对全服玩家说话，第二个参数为说的字符串
   - tell: 对某一玩家说话，第二个参数为玩家名称， 第二个参数为说的字符串
   - Execute: 执行某一命令，第二个参数为命令字符串
- ### 冷插件
   所有冷插件都在plugins文件夹中，包含一个栗子和一个基础插件，具体开发方式可以参考栗子，开发完以后，需要在pluginList.go中注册冷插件
- ### 自定义语法解释器
   所有语法解释器都在parsers中，包含一个默认语法解释器以及解释器列表，同样，开发完以后，需要在parserList.go中注册语法解释器

- ### `注意`
   所有插件和语法解释器都需要实现lib中对应的接口函数。