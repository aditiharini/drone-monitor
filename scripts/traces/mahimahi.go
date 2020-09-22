package trace

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type MahimahiTrace struct {
	Filename   string
	PacketSize int
}

func (m *MahimahiTrace) PrintBandwidth() {
	file, err := os.Open(m.Filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	lastCalcTime := -1
	numPackets := 0
	for scanner.Scan() {
		curTime, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(err)
		}

		if lastCalcTime == -1 {
			lastCalcTime = curTime
		}

		if curTime-lastCalcTime >= 1000 {
			numBits := float64(numPackets * m.PacketSize * 8)
			numMbits := numBits / 1000000.
			elapsed := float64(curTime-lastCalcTime) / 1000.
			mbps := numMbits / elapsed
			fmt.Println(mbps)

			numPackets = 0
			lastCalcTime = curTime
		}
		numPackets++
	}
}
