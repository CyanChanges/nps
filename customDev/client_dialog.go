package customDev

import (
	"github.com/astaxie/beego/logs"
	"github.com/lizongshen/gocommand"
	"log"
	"strings"
	"time"
)

func Dialog() {
	for {
		status := pppoeStatus()
		if status == "off" {
			// 直接拨号
			pppoeStart()

			for {
				time.Sleep(1 * time.Second)
				if pppoeStatus() == "on" {
					// 等待直到拨号成功
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
				break
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
