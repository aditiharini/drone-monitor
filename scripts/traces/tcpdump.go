package trace

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type TraceProcessor interface {
	ToMahiMahi()
}

type PcapProcessor struct {
	Filename           string
	Filter             string
	CurrentMahiFilenum int
	CurrentLossFilenum int
	MahimahiFiles      []string
	OutputDir          string
	FileDivisionTime   time.Duration
}

func (p *PcapProcessor) NewMahimahiTrace() *os.File {
	p.CurrentMahiFilenum++
	filename := fmt.Sprintf("uplink-%d.pps", p.CurrentMahiFilenum)
	p.MahimahiFiles = append(p.MahimahiFiles, filename)
	filePath := fmt.Sprintf("%s/%s", p.OutputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	return file
}

func (p *PcapProcessor) NewLossTrace() *os.File {
	p.CurrentLossFilenum++
	filename := fmt.Sprintf("uplink-%d.loss", p.CurrentLossFilenum)
	filePath := fmt.Sprintf("%s/%s", p.OutputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	return file
}

func (p PcapProcessor) LossAnalysis() {
	handle, err := pcap.OpenOffline(p.Filename)
	if err != nil {
		panic(err)
	}

	err = handle.SetBPFFilter(p.Filter)
	if err != nil {
		panic(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	// prevSeqno := uint32(0)
	maxSeqno := uint32(0)
	firstTime := -1
	numDropsEstimate := 0
	counter := 0
	for packet := range packetSource.Packets() {
		if counter < 100 {
			counter++
			continue
		}
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			arrivalTime := int(packet.Metadata().Timestamp.UnixNano() / 1000000)
			if firstTime == -1 {
				firstTime = arrivalTime
			}
			offset := arrivalTime - firstTime
			tcpPacket, _ := tcpLayer.(*layers.TCP)
			seqno := tcpPacket.Seq
			if seqno < maxSeqno {
				fmt.Println(offset, seqno, maxSeqno)
				numDropsEstimate++
			}
			if seqno > maxSeqno {
				maxSeqno = seqno
			}
			// prevSeqno = seqno
			counter++
		}
	}
	fmt.Printf("Num drops: %d, Num packets: %d\n", numDropsEstimate, counter)
}

func (p *PcapProcessor) ToLossTrace(granularity time.Duration) {
	lossCsvName := fmt.Sprintf("%s/uplink.csv", p.OutputDir)
	defer os.Remove(lossCsvName)

	tsharkCmd := exec.Command("bash", "-c", fmt.Sprintf("tshark -r %s -e frame.time_epoch -e tcp.analysis.lost_segment -e tcp.analysis.retransmission -T fields -E separator=, -E header=y > %s", p.Filename, lossCsvName))
	if err := tsharkCmd.Run(); err != nil {
		panic(err)
	}

	lossTraceFile := p.NewLossTrace()
	lossTraceWriter := csv.NewWriter(lossTraceFile)

	lossCsv, err := os.Open(lossCsvName)
	if err != nil {
		panic(err)
	}
	defer lossCsv.Close()

	lossCsvReader := csv.NewReader(lossCsv)
	lossCsvReader.Read()
	numPacketsInRange := 0
	numDropsInRange := 0
	rangeStart := time.Time{}
	firstTime := time.Time{}
	prevTime := time.Time{}
	for {
		row, err := lossCsvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		packetTime := row[0]
		splitTime := strings.Split(packetTime, ".")
		secs, err := strconv.ParseInt(splitTime[0], 10, 64)
		nanos, err := strconv.ParseInt(splitTime[1], 10, 64)
		packetTimeUnix := time.Unix(secs, nanos)

		if rangeStart.IsZero() {
			rangeStart = packetTimeUnix
		}

		if firstTime.IsZero() {
			firstTime = packetTimeUnix
		}

		if prevTime.IsZero() {
			prevTime = packetTimeUnix
		}

		if packetTimeUnix.Sub(prevTime) > p.FileDivisionTime {
			lossTraceWriter.Flush()
			lossTraceFile.Close()

			lossTraceFile = p.NewLossTrace()
			lossTraceWriter = csv.NewWriter(lossTraceFile)
			firstTime = packetTimeUnix
		}

		if packetTimeUnix.After(rangeStart.Add(granularity)) {
			if rangeStart.Before(firstTime) {
				rangeStart = firstTime
			}
			offsetMillis := rangeStart.Sub(firstTime) / time.Millisecond
			lossPercentage := float32(numDropsInRange) / float32(numPacketsInRange)
			row := []string{fmt.Sprintf("%d", offsetMillis), fmt.Sprintf("%f", lossPercentage)}
			lossTraceWriter.Write(row)

			rangeStart = packetTimeUnix
			numPacketsInRange = 0
			numDropsInRange = 0
		}

		// lostSegment := row[1]
		// if lostSegment == "1" {
		// 	numDropsInRange++
		// }

		retransmission := row[2]
		if retransmission == "1" {
			numDropsInRange++
		}
		numPacketsInRange++
		prevTime = packetTimeUnix
	}

	lossTraceWriter.Flush()
	lossTraceFile.Close()

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
		if offsetTime-prevTime > p.FileDivisionTime.Milliseconds() {
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
