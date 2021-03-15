package customDev

import (
	"fmt"
	"net"
)

func FindFreePort() (port int) {
	for i := ServerPortStart; i <= ServerPortEnd; i++ {
		if checkPort(i) > 0 {
			return i
		}
	}

	return
}

func checkPort(port int) int {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return 0
	}
	defer ln.Close()

	return ln.Addr().(*net.TCPAddr).Port
}
