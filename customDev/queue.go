package customDev

var (
	freePortsChan = make(chan int, 100000)
)

// 用队列来保障一定时间内不分配冲突的端口
func PopPort() int {
	if len(freePortsChan) <= 0 {
		restore()
	}

	return <-freePortsChan
}

// 重新填充
func restore() {
	for i := ServerPortStart; i <= ServerPortEnd; i++ {
		freePortsChan <- i
	}
}
