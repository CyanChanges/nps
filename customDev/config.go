package customDev

import (
	"time"
)

const (
	Debug     = true   // 是否开启调试模式
	ApiPort   = "8011" // Api端口
	ApiHost   = "127.0.0.1"
	PortStart = 20000 // 随机端口开始
	PortEnd   = 30000 // 随机端口结束
)

var (
	IpCreateTime time.Time // IP创建时间
)
