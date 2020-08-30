package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	interestmapsDir := "/home/pi/interestmaps/"
	if err := exec.Command("go", "run", interestmapsDir+"launch.go", interestmapsDir+"interestmaps-real.cfg", "drone", "0").Run(); err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)

	saturatrDir := "/home/pi/multisend/sender/"

	// Turn saturatr on and off
	for {
		saturatrCmd := exec.Command(saturatrDir+"saturatr", "ip", "eth1", "ip", "eth2", "3.91.1.79")
		if err := saturatrCmd.Start(); err != nil {
			fmt.Printf("Error starting saturatr %v\n", err)
		}
		time.Sleep(30 * time.Second)
		saturatrCmd.Process.Kill()
		time.Sleep(3 * time.Minute)
	}
}
