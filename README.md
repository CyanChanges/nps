# NPS (修改版)
[官方说明](https://github.com/ehang-io/nps/blob/master/README.md)|[中文文档](https://github.com/ehang-io/nps/blob/master/README_zh.md)

# nps 自定义开发说明

自定义开发主要是为了适应拨号VPS, 特别是内网拨号的主机而开发, 增加了自动获取服务器空闲端口的功能. 增加了定时重拨功能.

## 源码改动
- **所有扩展的功能全部放在 customDev 文件夹里**
- 文件夹内 client_ 前缀的是客户端调用的, server_ 前缀的是服务端调用的
- 服务端和客户端通过自定义的 web api 传递空闲端口
- 如果客户端连接分配的端口失败, 会重新从服务端获取新的端口
- 客户端每过N分钟重新拨号一次


## 注意
- 服务端自动分配 20000-30000 之间的端口, 请在防火墙上允许
- **目前只对 http 代理功能做过测试, 其他功能不保证能用.**
- 客户端的 conf 配置文件里的 server_port 已改为自动获取, conf 文件的可以忽略不管
- 客户端配置文件里 basic_username 和 basic_password 已改成了随机生成, 可以忽略不管

## 接口
- /api/freePort 用于客户端获取服务端可用端口
- /api/randHttpProxy/n 用于获取n个随机代理

## 配置文件
- 服务端安装之后他的配置文件在 /etc/nps/conf
- 客户端配置文件 npc.conf 里只保留了 http 隧道, 需要其他的要自己添加
_ 客户端和服务端的 disconnect_timeout 改成了5, 这样可以更快剔除断开的代理, 还可以更多测试后优化
- 配置文件其他参数说明:https://ehang-io.github.io/nps/#/server_config

## 服务器后台
-webApi http://127.0.0.1:8011
-后台 http://127.0.0.1:8080
-管理员:coco
-密码:Aio89:linux6