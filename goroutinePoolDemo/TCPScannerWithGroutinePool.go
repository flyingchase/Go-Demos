package main

import (
	"fmt"
	"net"
)

// 创建类似线程池避免 goroutine暴涨情况
// 限定 goroutine 创建的数量并复用
// channel传递 n 个数据分配到 m 个 goroutine
func main() {
	fmt.Println("TCP Scan start...")
	fmt.Println("the lists below is Open")
	n := make(chan int, 100)
	res := make(chan int, 100)
	// 创建 2w 个 goroutine
	for m := 0; m < 200; m++ {
		go worker(n, res)
	}
	// channel 代表任务数量
	for i := 1; i < 655; i++ {
		n <- i
	}
	close(n)
	for port := range res {
		fmt.Println(port)
	}
}

func worker(ports <-chan int, res chan<- int) {
	for port := range ports {
		address := fmt.Sprintf("192.168.1.109:%d", port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			continue
		}
		conn.Close()
		res <- port
	}
}
