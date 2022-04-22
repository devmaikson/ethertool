package ethernet

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/mdlayher/packet"
)

func GetNetworkInterface(interfaceName string) (*net.Interface, error) {
	// get a pointer to the network Interface (lookup by name)
	ifi, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("failed to find interface %q: %v", interfaceName, err)
	}

	return ifi, nil
}

func OpenConnection(ifi *net.Interface, etherType int) (net.PacketConn, error) {
	//open raw socket using packet library by mdlayher
	conn, err := packet.Listen(ifi, packet.Raw, etherType, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open socket %v", err)
	}

	return conn, err
}

// sendMessages continuously sends a message over a connection at regular intervals,
// sourced from specified hardware address.
func SendMessages(c net.PacketConn, destAddr net.HardwareAddr, f *Frame, interval time.Duration) {

	binary, err := f.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to marshal ethernet frame: %v", err)
	}

	// Required by Linux, even though the Ethernet frame has a destination.
	// Unused by BSD.
	addr := &packet.Addr{
		HardwareAddr: destAddr,
	}

	// Send message forever.
	t := time.NewTicker(interval * time.Second)
	for range t.C {
		if _, err := c.WriteTo(binary, addr); err != nil {
			log.Fatalf("failed to send message: %v", err)
		}
		print(".")
	}
}

// receiveMessages continuously receives messages over a connection. The messages
// may be up to the interface's MTU in size.
func ReceiveMessages(c net.PacketConn, mtu int) {
	var f Frame
	b := make([]byte, mtu)

	// Keep receiving messages forever.
	for {
		n, addr, err := c.ReadFrom(b)
		if err != nil {
			log.Fatalf("failed to receive message: %v", err)
		}

		// Unpack Ethernet II frame into Go representation.
		if err := (&f).UnmarshalBinary(b[:n]); err != nil {
			log.Fatalf("failed to unmarshal ethernet frame: %v", err)
		}

		// Display source of message and message itself.
		log.Printf("\n[%s] %s", addr.String(), string(f.Payload))
	}
}

func DefaultMsg(defaultMsg string) (msg string, err error) {
	// Default message to system's hostname if empty.
	if defaultMsg == "" {
		newMsg, err := os.Hostname()
		if err != nil {
			return newMsg, fmt.Errorf("failed to retrieve hostname: %v", err)
		}
		return newMsg, nil
	}

	return defaultMsg, nil
}
