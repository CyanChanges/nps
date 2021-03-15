package customDev

import (
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
		// 把新IP给代理池
		//submit()

		for {
			// 等到IP过期后重新拨号
			if IpExpired() {
				break
			}
			if Debug {
				println("Proxy is running")
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func pppoeStart() (result bool) {
	_, _, err := gocommand.NewCommand().Exec("pppoe-start")
	if err != nil {
		println("Get err when start ppoe: ", err)
		return
	}
	println("pppoe start")
	IpCreateTime = time.Now()
	return true
}

func pppoeStop() (result bool) {
	_, _, err := gocommand.NewCommand().Exec("pppoe-stop")
	if err != nil {
		println("Get err when stop ppoe: ", err)
		return
	}
	println("pppoe stop")
	return true
}

func pppoeStatus() (status string) {
	_, out, err := gocommand.NewCommand().Exec("pppoe-status")
	if err != nil {
		log.Panic(err)
	}

	if Debug {
		println(out)
	}

	if strings.Contains(out, "Link is up and running") {
		return "on"
	} else if strings.Contains(out, "Link is down") {
		return "off"
	}

	return
}

// IP过期时间检测
func IpExpired() (status bool) {
	if time.Now().Sub(IpCreateTime) >= time.Duration(60)*time.Second {
		status = true
	}
	return
}
