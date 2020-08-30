package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i <= 1024; i++ {
		wg.Add(1)
		// run the scans as a goroutine
		go func(j int) {
			defer wg.Done()
			address := fmt.Sprintf("scanme.nmap.org:%d", j)
			//fmt.Println(address)

			conn, err := net.DialTimeout("tcp", address, 5*time.Second)
			if err != nil {
				// port is closed or filtered
				return
			}
			conn.Close()
			fmt.Printf("%d open\n", j)
		}(i)
	}
	wg.Wait()
}
