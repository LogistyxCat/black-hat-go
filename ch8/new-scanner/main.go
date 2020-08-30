package main

/*
TO DO:
	[ ] Remove the network and transport layer logic and source checks from the capture() func
		 - This should be rolled into the BPF filter to ensure only those packets are captured
	[ ] Replace the sequestionl logic of port scanning with a concurrent alternative
	[ ] Rather than limiting the code to a single target, allow the user to supply a list of IPs or network blocks
*/

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// We need to check that service is listening on the port.
// These flag combinations are common for real open ports, so we can calculate the filter with these.
// 		ACK && FIN:	00010001	0x11
// 		ACK:		00010000	0x10
// 		ACK && PSH:	00011000	0x18
// BPF Filter:	tcp[13] == 0x11 or tcp[13] == 0x10 or tcp[13] == 0x18

var (
	snaplen  = int32(1600)
	promisc  = true
	timeout  = pcap.BlockForever
	filter   = "tcp[13] == 0x11 or tcp[13] == 0x10 or tcp[13] == 0x18"
	devFound = false
	results  = make(map[string]int)
)

func capture(iface, target string) {
	handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	if err != nil {
		log.Panicln(err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter(filter); err != nil {
		log.Panicln(err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	fmt.Println("[*] Capturing packets...")
	for packet := range source.Packets() {
		networkLayer := packet.NetworkLayer()
		if networkLayer == nil {
			continue
		}
		transportLayer := packet.TransportLayer()
		if transportLayer == nil {
			continue
		}

		srcHost := networkLayer.NetworkFlow().Src().String()
		srcPort := transportLayer.TransportFlow().Src().String()

		if srcHost != target {
			continue
		}

		// Increment confidence value
		results[srcPort]++
	}
}

func explode(portString string) ([]string, error) {
	ret := make([]string, 0)

	ports := strings.Split(portString, ",")
	for _, port := range ports {
		port := strings.TrimSpace(port)
		ret = append(ret, port)
	}

	return ret, nil
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalln("Usage:", os.Args[0], "<capture_iface> <target_ip> <port1,port2,port3>")
	}

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}

	iface := os.Args[1]
	for _, device := range devices {
		if device.Name == iface {
			devFound = true
		}
	}
	if !devFound {
		log.Panicf("Device named '%s' does not exist\n", iface)
	}

	ip := os.Args[2]
	go capture(iface, ip)
	time.Sleep(1 * time.Second)

	ports, err := explode(os.Args[3])
	if err != nil {
		log.Panicln(err)
	}

	for _, port := range ports {
		target := fmt.Sprintf("%s:%s", ip, port)
		fmt.Println("[*] Trying", target)
		c, err := net.DialTimeout("tcp", ip, 1000*time.Millisecond)
		if err != nil {
			continue
		}
		c.Close()
	}
	time.Sleep(2 * time.Second)

	for port, confidence := range results {
		if confidence >= 1 {
			fmt.Printf("Port %s open (confidence: %d)\n", port, confidence)
		}
	}
}
