package client

import "time"

const Debug = true                                 // 是否开启调试模式
const ProxySubmitUrl = "http://localhost:8081/add" // 代理提交地址

var (
	IpCreateTime time.Time // IP创建时间
)
