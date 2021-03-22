package customDev

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lizongshen/gocommand"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func Dialog() {
	for {
		status := pppoeStatus()
	rePppoeImmediately:
		if status == "off" {
			// 直接拨号
			pppoeStart()

			for {
				time.Sleep(1 * time.Second)
				if pppoeStatus() == "on" {
					// 等待直到拨号成功
					ClientDisInternet = false
					break
				}
			}
		} else if status == "on" {
			// 重新拨号
			pppoeStop()
			for {
				time.Sleep(1 * time.Second)
				if pppoeStatus() == "off" {
					// 等待直到断开拨号
					break
				}
			}
			pppoeStart()
			for {
				time.Sleep(1 * time.Second)
				if pppoeStatus() == "on" {
					// 等待直到拨号成功
					ClientDisInternet = false
					break
				}
			}
		} else {
			time.Sleep(1 * time.Second)
			continue
		}
		time.Sleep(5 * time.Second)

		for {
			// 等到IP过期后重新拨号
			if IpExpired() {
				tellServerDel()

				// 马上进入重新拨号流程
				status = "on"
				goto rePppoeImmediately
			}
			//logs.Error("Proxy is running")
			time.Sleep(1 * time.Second)
		}
	}
}

func pppoeStart() (result bool) {
	_, _, err := gocommand.NewCommand().Exec("pppoe-start")
	if err != nil {
		logs.Error("Get err when start pppoe: ", err)
		return
	}
	logs.Info("pppoe start")
	ClientIpCreateTime = time.Now()
	return true
}

func pppoeStop() (result bool) {
	_, _, err := gocommand.NewCommand().Exec("pppoe-stop")
	if err != nil {
		logs.Error("Get err when start pppoe: ", err)
		return
	}
	logs.Info("pppoe stop")
	return true
}

func pppoeStatus() (status string) {
	_, out, err := gocommand.NewCommand().Exec("pppoe-status")
	if err != nil {
		log.Panic(err)
	}

	//println(out)

	if strings.Contains(out, "Link is up and running") {
		return "on"
	} else if strings.Contains(out, "Link is down") {
		return "off"
	} else {
		logs.Error("pppoe-status returns nothing")
	}

	return
}

// IP过期时间检测
func IpExpired() (status bool) {
	if time.Now().Sub(ClientIpCreateTime) >= time.Duration(ClientPppoeExpiry)*time.Second {
		status = true
	}
	return
}

// 通知服务端删除自身, 然后等待客户端隧道退出
func tellServerDel() {
	ClientDisInternet = true

	logs.Notice("告诉主机我要换IP了")
	c := &http.Client{
		Timeout: 30 * time.Second,
	}
	r, err := c.Get(fmt.Sprintf("http://%s:%s/api/delClient", ApiHost, ApiPort))
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

	logs.Notice("通知服务器删除自身: %s\n", body)

	for {
		time.Sleep(time.Second)
		if ClientGotDelFlag {
			ClientGotDelFlag = false
			break
		}
	}
}
