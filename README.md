# NPS (修改版)
[官方说明](https://github.com/ehang-io/nps/blob/master/README.md)|[中文文档](https://github.com/ehang-io/nps/blob/master/README_zh.md)

# nps 自定义开发说明

自定义开发主要是为了适应拨号VPS, 特别是内网拨号的主机而开发, 增加了自动获取服务器空闲端口的功能. 增加了定时拨号功能.

## 源码改动
- **所有扩展的功能全部放在 customDev 文件夹里**
- 文件夹内 client_ 前缀的是客户端调用的, server_ 前缀的是服务端调用的
- 服务端和客户端通过自定义的 web api 传递空闲端口
- 如果客户端连接分配的端口失败, 会重新从服务端获取新的端口
- 客户端每过1分钟重新拨号一次


## 注意
- **目前只对 http 代理功能做过测试, 其他功能不保证能用.**
- 客户端的 conf 配置文件里的 server_port 参数就相当与作废了
