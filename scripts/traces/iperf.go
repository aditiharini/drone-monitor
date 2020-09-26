package trace

import (
	"bufio"
	"fmt"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type TraceProcessor interface {
	ToMahiMahi()
}

type PcapProcessor struct {
	Filename       string
	Filter         string
	CurrentFilenum int
	MahimahiFiles  []string
	OutputDir      string
}

func (p *PcapProcessor) NewMahimahiTrace() *os.File {
	p.CurrentFilenum++
	filename := fmt.Sprintf("uplink-%d.pps", p.CurrentFilenum)
	p.MahimahiFiles = append(p.MahimahiFiles, filename)
	filePath := fmt.Sprintf("%s/%s", p.OutputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	return file
}

func (p *PcapProcessor) ToMahiMahi() {
	handle, err := pcap.OpenOffline(p.Filename)
	if err != nil {
		panic(err)
	}

	err = handle.SetBPFFilter(p.Filter)
	if err != nil {
		panic(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	firstTime := int64(-1)
	prevTime := int64(-1)
	mahimahiFile := p.NewMahimahiTrace()
	mahimahiWriter := bufio.NewWriter(mahimahiFile)
	for packet := range packetSource.Packets() {
		if firstTime == -1 {
			firstTime = packet.Metadata().Timestamp.UnixNano()
		}
		offsetTime := (packet.Metadata().Timestamp.UnixNano() - firstTime) / int64(1000000)
		if prevTime == -1 {
			prevTime = offsetTime
		}
		if offsetTime-prevTime > 50000 {
			mahimahiWriter.Flush()
			if err := mahimahiFile.Close(); err != nil {
				panic(err)
			}

			firstTime = packet.Metadata().Timestamp.UnixNano()
			offsetTime = 0
			mahimahiFile = p.NewMahimahiTrace()
			mahimahiWriter = bufio.NewWriter(mahimahiFile)
		}
		mahimahiWriter.WriteString(fmt.Sprintf("%d\n", offsetTime))
		prevTime = offsetTime
	}
}
