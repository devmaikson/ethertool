# ethertool
Go based tool to send and receive Ethernet frames

## Usage

```
sudo ./ethertool -h
Usage of ./ethertool:
  -cEthType int
    	customer VLAN Ethertype (default 0x8100)
  -cpbit int
    	customer VLAN Pbit (default 7)
  -cvid int
    	customer VLAN-iD (default 4)
  -d string
    	destination mac-address to use to send messages (default: broadcast)
  -i string
    	network interface to use to send and receive messages
  -m string
    	message to be sent (default: system's hostname)
  -sEthType int
    	Outer VLAN Ethertype (default 0x8100)
  -spbit int
    	s-tag VLAN Pbit (default 3)
  -svid int
    	s-tag VLAN-iD (default 3200)
  -t int
    	time interval for sending messages in seconds) (default 2)
```

## Example

```
sudo ./ethertool -i enp0s8
Start sending packets every: 2s
..
```
on another terminal:
```
sudo tcpdump -i enp0s8 -vvv -w capture.pcap
```
