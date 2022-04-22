package main

import (
	"flag"
	"fmt"
	"log"
	"maikson/ethertool/internal/pkg/ethernet"
	"net"
	"time"
)

// Ethertype of the Packet
// https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml
const etherType = 0xcccc

func main() {
	var (
		ifaceFlag     = flag.String("i", "", "network interface to use to send and receive messages")
		destMacFlag   = flag.String("d", "", "destination mac-address to use to send messages (default: broadcast)")
		msgFlag       = flag.String("m", "", "message to be sent (default: system's hostname)")
		timeFlag      = flag.Int("t", 2, "time interval for sending messages in seconds)")
		sVlanFlag     = flag.Int("svid", 3200, "s-tag VLAN-iD")
		cVlanFlag     = flag.Int("cvid", 4, "customer VLAN-iD")
		cVlanPbitFlag = flag.Int("cpbit", 7, "customer VLAN Pbit")
		sVlanPbitFlag = flag.Int("spbit", 3, "s-tag VLAN Pbit")
		sEthTypeFlag  = flag.Int("sEthType", 0x8100, "Outer VLAN Ethertype (default 0x8100)")
		cEthTypeFlag  = flag.Int("cEthType", 0x8100, "customer VLAN Ethertype (default 0x8100)")
	)
	flag.Parse()

	iface, err := ethernet.GetNetworkInterface(*ifaceFlag)
	if err != nil {
		log.Fatalf("error getting network interface: %v", err)
	}

	conn, err := ethernet.OpenConnection(iface, etherType)
	if err != nil {
		log.Fatalf("error creating a connection: %v", err)
	}

	// default msg is hostname if not provided as a flag
	msg, err := ethernet.DefaultMsg(*msgFlag)
	if err != nil {
		log.Fatalf("error reading hostname try setting message with -m \"your-test\" %v", err)
	}

	// destination mac-address is set to broadcast by default
	dstMac := ethernet.Broadcast
	// if destination mac is provided as a flag use it
	if *destMacFlag != "" {
		mac, err := net.ParseMAC(*destMacFlag)
		if err != nil {
			log.Fatalf("error parsing mac-address: %v", err)
		}
		dstMac = mac
	}

	ethernet.EtherTypeServiceVLAN = ethernet.EtherType(uint16(*sEthTypeFlag))
	ethernet.EtherTypeVLAN = ethernet.EtherType(uint16(*cEthTypeFlag))

	frame := &ethernet.Frame{
		Destination: dstMac,
		Source:      iface.HardwareAddr,
		ServiceVLAN: &ethernet.VLAN{Priority: ethernet.Priority(uint8(*sVlanPbitFlag)), DropEligible: false, ID: uint16(*sVlanFlag)},
		VLAN:        &ethernet.VLAN{Priority: ethernet.Priority(uint8(*cVlanPbitFlag)), DropEligible: false, ID: uint16(*cVlanFlag)},
		EtherType:   etherType,
		Payload:     []byte(msg),
	}

	fmt.Printf("Start sending packets every: %vs\n", *timeFlag)
	// Send messages in one goroutine, receive messages in another.
	go ethernet.SendMessages(conn, dstMac, frame, time.Duration(*timeFlag))
	go ethernet.ReceiveMessages(conn, iface.MTU)

	// Block forever.
	select {}
}
