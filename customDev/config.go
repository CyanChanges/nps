package customDev

import (
	"bufio"
	"ehang.io/nps/lib/common"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

const (
	ApiPort           = "8011" // Api端口
	ServerPortStart   = 20000  // 服务端随机端口开始, 请注意防火墙
	ServerPortEnd     = 30000  // 服务端随机端口结束, 请注意防火墙
	serverPppoeExpiry = 120    // 客户端通过api获取, 拨号间隔(秒)
)

var (
	ApiHost = getApiHost()
	//ApiHost            = "127.0.0.1"
	ClientIpCreateTime time.Time          // IP创建时间
	ClientPppoeExpiry  = getPppoeExpiry() // 客户端获取
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

func getPppoeExpiry() (port int) {
	r, err := http.Get(fmt.Sprintf("http://%s:%s/api/pppoeExpiry", ApiHost, ApiPort))
	if err != nil {
		logs.Error("代理池的API无法访问: %s\n", err)
		return
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil || r.StatusCode != 200 {
		logs.Error("代理池的API无法访问: %s %d\n", err, r.StatusCode)
		return
	}

	// 判断是否数字
	expiryInt, err := strconv.Atoi(string(body))
	if err != nil {
		return serverPppoeExpiry // 如果服务端没有就使用本地默认的
	}

	return expiryInt
}
