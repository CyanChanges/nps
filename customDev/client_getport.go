package customDev

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetPort() (port string) {
	r, err := http.Get(fmt.Sprintf("http://%s:%s/api/freePort", ApiHost, ApiPort))
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

	port = string(body)

	// 判断是否数字
	_, err = strconv.Atoi(port)
	if err != nil {
		return
	}

	return port
}
