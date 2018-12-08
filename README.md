# MCDaemon-go

## 用golang实现的Minecraft进程管理程序

### 在windows10以及centos7中运行成功
-----
## 开始使用
- ### 下载最新的 release
- ### 修改配置文件MCD_conig.ini
  1. 修改server_name为服务端文件名
  2. 修改server_path为服务端文件所在文件夹，注意路径！
- ### 根据需求更改配置文件参数
- ### 运行start(linux/unix)或者start.exe（windows）
-----
## 插件编写

- ### 热插件
   1. 将热插件的执行文件放入hotPlugins文件夹中
   2. 在MCD_conig.ini的[plugins]域中注册热插件
   3. 通过：
        ```go
            args := os.Args
	        args = args[1:] 
        ```
      获取插件的命令的参数
   4. 返回调用的MCD命令api以及参数
- ### 冷插件
- ### 自定义语法解析器
