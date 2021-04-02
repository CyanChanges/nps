package customDev

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func GetPort() (port string) {
	c := &http.Client{
		Timeout: 30 * time.Second,
	}
	r, err := c.Get(fmt.Sprintf("http://%s:%s/api/freePort", ApiHost, ApiPort))

	if err != nil {
		logs.Error("代理池的API无法访问: %s", err)
		return
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil || r.StatusCode != 200 {
		logs.Error("代理池的API无法访问: %s %d", err, r.StatusCode)
		return
	}

	port = string(body)

	// 判断是否数字
	_, err = strconv.Atoi(port)
	if err != nil {
		return
	}

	return port
}
