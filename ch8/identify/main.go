package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket/pcap"
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}

	for _, device := range devices {
		fmt.Println(device.Name)
		for _, address := range device.Addresses {
			// Check if IP is v4 or v6
			if address.IP.To4() != nil {
				fmt.Printf("\tIPv4:\t %s\n", address.IP)
			} else {
				fmt.Printf("\tIPv6:\t %s\n", address.IP)
			}
			fmt.Printf("\tNetmask: %s\n", address.Netmask)
		}
	}
}
