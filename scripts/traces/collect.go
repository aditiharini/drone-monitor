package trace

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func ProcessTrace(pcapfile string, filter string) {
	handle, err := pcap.OpenOffline(pcapfile)
	if err != nil {
		panic(err)
	}

	err = handle.SetBPFFilter(filter)
	if err != nil {
		panic(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	firstTime := int64(-1)
	for packet := range packetSource.Packets() {
		if firstTime == -1 {
			firstTime = packet.Metadata().Timestamp.UnixNano()
		}
		fmt.Println((packet.Metadata().Timestamp.UnixNano() - firstTime) / int64(1000000))
	}
}
