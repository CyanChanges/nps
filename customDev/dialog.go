package client

import (
	"github.com/lizongshen/gocommand"
	"log"
	"strings"
	"time"
)

func PpoeStart() (result bool) {
	_, _, err := gocommand.NewCommand().Exec("pppoe-start")
	if err != nil {
		println("Get err when start ppoe: ", err)
		return
	}
	println("pppoe start")
	return true
}

func PpoeStop() (result bool) {
	_, _, err := gocommand.NewCommand().Exec("pppoe-stop")
	if err != nil {
		println("Get err when stop ppoe: ", err)
		return
	}
	println("pppoe stop")
	return true
}

func PppoeStatus() (status string) {
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
