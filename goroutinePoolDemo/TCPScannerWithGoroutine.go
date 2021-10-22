package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	fmt.Println("TCP端口扫描启动...")
	start := time.Now()
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			address := fmt.Sprintf("192.168.1.109:%d", i)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				fmt.Printf(" %s Closed!\n", address)
				return
			}
			conn.Close()
			fmt.Printf(" %s Open\n", address)
		}(i)
	}
	wg.Wait()
	elasped := time.Since(start)/1e9
	fmt.Printf("\n 经过 %d 秒\n", elasped)
}
