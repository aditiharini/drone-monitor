package main

import "os/exec"

func main() {
	exec.Command("killall", "interestmaps").Run()
	exec.Command("killall", "drone-monitor").Run()
	exec.Command("killall", "/usr/local/bin/serve").Run()
	exec.Command("killall", "iperf3").Run()
	exec.Command("killall", "tcpdump").Run()
}
