package customDev

import (
	"bufio"
	"ehang.io/nps/lib/common"
	"os"
	"regexp"
	"time"
)

const (
	ApiPort         = "8011" // Api端口
	ServerPortStart = 20000  // 服务端随机端口开始, 请注意防火墙
	ServerPortEnd   = 30000  // 服务端随机端口结束, 请注意防火墙
)

var (
	ApiHost            = getApiHost()
	ClientIpCreateTime time.Time // IP创建时间
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
