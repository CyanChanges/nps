package customDev

import (
	"ehang.io/nps/gocommand"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/parnurzeal/gorequest"
	"strings"
	"syscall"
	"time"
)

func Adsl() {
	for {
		status := pppoeStatus()
	rePppoeImmediately:
		if status == "off" {
			// 直接拨号
			if pppoeStart() {
				for i := 1; i <= 8; i++ {
					time.Sleep(1 * time.Second)
					if pppoeStatus() == "on" {
						// 等待直到拨号成功
						break
					}
				}
			}

		} else if status == "on" {
			time.Sleep(2 * time.Second)
			// 重新拨号
			if pppoeStop() {
				for i := 1; i <= 8; i++ {
					time.Sleep(1 * time.Second)
					if pppoeStatus() == "off" {
						// 等待直到断开拨号
						break
					}
				}
			}

			if pppoeStart() {
				for i := 1; i <= 8; i++ {
					time.Sleep(1 * time.Second)
					if pppoeStatus() == "on" {
						// 等待直到拨号成功
						break
					}
				}
			}

		} else {
			// Pppoe-status 查询失败的情况
			time.Sleep(5 * time.Second)
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
	_, success := cmd("/usr/sbin/pppoe-start")

	if success == false {
		time.Sleep(time.Second)
		return
	}

	logs.Info("pppoe start")
	ClientIpCreateTime = time.Now()
	return true
}

func pppoeStop() (result bool) {
	_, success := cmd("/usr/sbin/pppoe-stop")

	if success == false {
		time.Sleep(time.Second)
		return
	}

	logs.Info("pppoe stop")
	return true

}

func pppoeStatus() (status string) {
	out, success := cmd("/usr/sbin/pppoe-status")
	if success == false {
		logs.Error("pppoe-status failed")
		return
	}

	if strings.Contains(out, "Link is up and running") {
		return "on"
	} else if strings.Contains(out, "Link is down") {
		return "off"
	}

	logs.Error("pppoe-status return unexpect: ", out)
	if strings.Contains(out, "ppp2 is down") {
		pppoeStop()
		time.Sleep(time.Minute)
	}
	return
}

// IP过期时间检测
func IpExpired() (status bool) {
	if time.Now().Sub(ClientIpCreateTime) >= time.Duration(ClientPppoeExpiry)*time.Second {
		return true
	}
	return
}

// 通知服务端删除自身, 然后等待客户端隧道退出
func tellServerDel() {
	defer func() {
		if err1 := recover(); err1 != nil {
			logs.Error("Command 发生严重错误", err1)
		}
	}()

	logs.Notice("告诉主机我要换IP了")

	request := gorequest.New()

	for i := 1; i <= 5; i++ {
		resp, body, errs := request.Get(fmt.Sprintf("http://%s:%s/api/delClient", ApiHost, ApiPort)).
			Timeout(30 * time.Second).
			End()

		if err := ErrAndStatus(errs, resp); err != nil {
			logs.Error("代理池的API无法访问: %s", err)
			time.Sleep(5 * time.Second)
			continue
		}

		logs.Notice("通知服务器删除自身: %s", body)
		time.Sleep(time.Second)
		break
	}
}

func cmd(command string) (result string, success bool) {
	defer func() {
		if err := recover(); err != nil {
			logs.Error("Command 发生严重错误", err)
		}
	}()

	var cmd, out, err = gocommand.NewCommand().Exec(command)

	if cmd != nil {
		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err == nil {
			errKill := syscall.Kill(-pgid, 15) // note the minus sign
			if errKill != nil {
				logs.Error("Kill 发生错误", errKill)
			}
		}

		_ = cmd.Wait()
	}

	if err != nil {
		logs.Error("执行命令发生错误 "+command+": ", err)
		return
	}

	return out, true
}
