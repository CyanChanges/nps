# NPS (修改版)
[官方说明](https://github.com/ehang-io/nps/blob/master/README.md)|[中文文档](https://github.com/ehang-io/nps/blob/master/README_zh.md)

# 使用说明
服务端:
- 1.将 nps 和 conf 传到服务器
- 2.使用 sudo ./nps install 安装服务端
- 3.然后使用 sudo nps start 启动服务端

客户端:
- 1.将 client, npc, conf 传到拨号机器
- 2.npc.conf 里的 vkey 请配置一个唯一标识符, server_addr 配置远程服务器
- 3.在服务端的管理系统手动添加该客户端配置, 然后给该客户端添加对应的隧道方式
- 4.最后运行 client 即可

## 源码改动
- **所有扩展的功能全部放在 customDev 文件夹里**
- 文件夹内 client_和npc_ 前缀的是客户端调用的, server_ 前缀的是服务端调用的
- 服务端和客户端通过自定义的 web api 传递空闲端口
- 客户端每过N分钟重新拨号一次

## 注意
- 请在服务端防火墙上允许 7999, 8002, 8080 以及手动给代理配置的端口,例如 10000-10100
- **目前只对 http 代理方式做过测试**
- **后面会开发批量自动配置功能**

## 接口
- /api/freePort 用于客户端获取服务端可用端口
- /api/randHttpProxy/n 用于获取n个随机代理

## 配置文件
- 服务端安装之后他的配置文件在 /etc/nps/conf
- 客户端配置文件则放到 npc 和 client 的目录下即可
- 配置文件其他参数说明:https://ehang-io.github.io/nps/#/server_config

## 服务器后台,请在 nps.conf 里自行修改密码
-webApi http://127.0.0.1:8002
-管理后台 http://127.0.0.1:8080
-管理员:admin
-密码:admin123