package main

import (
	"flag"

	trace "github.com/aditiharini/drone-monitor/scripts/traces"
)

func main() {
	tracefile := flag.String("trace", "", "trace to be processed")
	isTcp := flag.Bool("tcp", false, "is this a tcp trace")
	flag.Parse()

	proto := "udp"
	if *isTcp {
		proto = "tcp"
	}
	processor := trace.PcapProcessor{Filename: *tracefile, Filter: proto, CurrentFilenum: 0}
	processor.ToMahiMahi()
	for _, file := range processor.MahimahiFiles {
		mmTrace := trace.MahimahiTrace{Filename: file, PacketSize: 1500}
		mmTrace.PrintBandwidth(false)
	}
}
