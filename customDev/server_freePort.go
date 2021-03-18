package customDev

import (
	"fmt"
	"net"
)

func FindFreePort() (port int) {
	for {
		port = CheckPort(PopPort())
		if port > 0 {
			return port
		}
	}
}

// 检查端口是否被占用
func CheckPort(port int) int {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return 0
	}
	defer ln.Close()

	return ln.Addr().(*net.TCPAddr).Port
}
