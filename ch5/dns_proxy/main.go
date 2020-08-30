package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/miekg/dns"
)

func parse(filename string) (map[string]string, error) {
	records := make(map[string]string)
	fh, err := os.Open(filename)
	if err != nil {
		return records, err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 2)
		if len(parts) < 2 {
			return records, fmt.Errorf("[!] %s is not a valid line", line)
		}
		records[parts[0]] = parts[1]
	}
	log.Println("[*] Records set to:")
	for k, v := range records {
		fmt.Printf("\t%s -> %s\n", k, v)
	}
	return records, scanner.Err()
}

func main() {
	var recordLock sync.RWMutex

	records, err := parse("proxy.conf")
	if err != nil {
		panic(err)
	}

	dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		if len(req.Question) < 1 {
			dns.HandleFailed(w, req)
			return
		}
		fqdn := req.Question[0].Name
		parts := strings.Split(fqdn, ".")
		if len(parts) < 1 {
			fqdn = strings.Join(parts[len(parts)-2:], ".")
		}

		// we need to strip any trailing periods (.)
		if fqdn[len(fqdn)-1:] == "." {
			fmt.Println("[Debug] Request ends in period, stripping...")
			fqdn = strings.TrimRight(fqdn, ".")
		}
		recordLock.RLock()
		match, ok := records[fqdn]
		recordLock.RUnlock()
		if !ok {
			fmt.Println("[!] No records found for", fqdn)
			dns.HandleFailed(w, req)
			return
		}
		fmt.Println("[*] Record found for", fqdn, "->", match)
		resp, err := dns.Exchange(req, match)
		if err != nil {
			fmt.Println("[!] Could not forward request to upstream DNS server")
			dns.HandleFailed(w, req)
		}
		if err := w.WriteMsg(resp); err != nil {
			fmt.Println("[!] Could not write response to client")
			dns.HandleFailed(w, req)
			return
		}
	})

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGUSR1)

		for sig := range sigs {
			switch sig {
			case syscall.SIGUSR1:
				recordLock.Lock()
				parse("proxy.conf")
				recordLock.Unlock()
			}
		}
	}()

	log.Fatal(dns.ListenAndServe(":553", "udp", nil))
}
