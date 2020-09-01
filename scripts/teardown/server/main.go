package main

import "os/exec"

func main() {
	exec.Command("killall", "saturatr").Run()
	exec.Command("killall", "interestmaps").Run()
	exec.Command("killall", "drone-monitor").Run()
	exec.Command("killall", "/usr/local/bin/serve").Run()
}
