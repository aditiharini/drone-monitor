package main

import (
	"os/exec"
	"time"
)

func main() {
	monitorDir := "/home/ubuntu/data/drone-monitor/"
	if err := exec.Command(monitorDir + "drone-monitor").Run(); err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	interestmapsDir := "/home/ubuntu/interestmaps/"
	if err := exec.Command(interestmapsDir+"interestmaps", interestmapsDir+"interestmaps-real.cfg"); err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	saturatrDir := "/home/ubuntu/data/multisend/sender/"
	if err := exec.Command(saturatrDir+"saturatr", "real").Run(); err != nil {
		panic(err)
	}
}
