package main

import (
	"fmt"
	"net"
	"sort"
	"time"
)

// Worker routine to test connection to ports fed into ports channel.
// Exports port number to results channel if successful, otherwise '0'.
func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.DialTimeout("tcp", address, 5*time.Second)

		if err != nil {
			results <- 0
			continue
		}

		// idea: implement verbosity flag and include this line
		//fmt.Println(fmt.Sprintf("[*] Found port %d", p))

		conn.Close()
		results <- p
	}
}

func main() {
	ports := make(chan int, 75) // idea: implement method to specify number of workers
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	// concurrently feed port list into ports channel for processing
	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}
	close(ports)
	close(results)

	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("[+] Port %d open\n", port)
	}

	fmt.Println("Scan complete")
}
