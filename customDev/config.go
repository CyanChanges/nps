package customDev

import (
	"bufio"
	"ehang.io/nps/lib/common"
	"os"
	"regexp"
	"time"
)

const (
	ApiPort           = "8011" // Api端口
	ServerPortStart   = 20000  // 服务端随机端口开始, 请注意防火墙
	ServerPortEnd     = 30000  // 服务端随机端口结束, 请注意防火墙
	ClientPppoeExpiry = 120    // 拨号间隔(秒)
)

var (
	ApiHost = getApiHost()
	//ApiHost            = "127.0.0.1"
	ClientIpCreateTime time.Time // IP创建时间
	//ClientDisInternet  bool      // 客户端断开互联网
	//ClientGotDelFlag   bool      // 客户端删除后, tcp会收到服务器已断开的信号
)

// 从配置文件中获取服务端地址
func getApiHost() (host string) {
	serverAddr := readLine(common.GetConfigPath(), 2)
	if serverAddr == "" {
		return
	}
	println(serverAddr)
	re := regexp.MustCompile(`server_addr=(((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}):\d{3,5}`)
	match := re.FindStringSubmatch(serverAddr)
	return match[1]
}

// 读取文件指定行
func readLine(filePath string, lineNumber int) string {
	file, err := os.Open(filePath)

	if err != nil {
		return ""
	}

	fileScanner := bufio.NewScanner(file)
	lineCount := 1
	for fileScanner.Scan() {
		if lineCount == lineNumber {
			return fileScanner.Text()
		}
		lineCount++
	}
	defer file.Close()
	return ""
}
