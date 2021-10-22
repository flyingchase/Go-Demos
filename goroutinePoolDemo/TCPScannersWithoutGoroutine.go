package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("TCP端口扫描启动...")
	start := time.Now()
	for i := 1; i <= 20; i++ {
		address := fmt.Sprintf("192.168.1.109:%d", i)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Printf("%s Closed!\n", address)
			continue
		}
		conn.Close()
		fmt.Printf(" %s Open!\n", address)
	}
	elasped := time.Since(start)/1e9
	fmt.Printf("\n 经过了%d 秒\n", elasped)
}
