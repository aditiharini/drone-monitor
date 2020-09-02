package main

import (
	"fmt"
	"os/exec"
	"time"
	"bufio"
)
func run(cmd *exec.Cmd, tag string, printStdout bool, printStderr bool) {
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	if printStdout {
		go func() {
			buf := bufio.NewReader(stdout)
			for {
				line, _, err := buf.ReadLine()
				if err != nil {
					break
				}
				fmt.Println(tag, string(line))
			}
		}()
	}

	if printStderr {
		go func() {
			buf := bufio.NewReader(stderr)
			for {
				line, _, err := buf.ReadLine()
				if err != nil {
					break
				}
				fmt.Println(tag, string(line))
			}
		}()
	}
}

func setupDevices() {
	if err := exec.Command("bash", "-c", "sudo ifconfig eth2 down").Run(); err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	if err := exec.Command("bash", "-c", "sudo ifconfig eth2 hw ether 0c:5b:8f:27:9a:65").Run(); err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	if err := exec.Command("bash", "-c", "sudo ifconfig eth2 up").Run(); err != nil {
		panic(err)
	}
	// time.Sleep(1 * time.Second)
	// if err := exec.Command("bash", "-c", "sudo ifconfig wlan0 down").Run(); err != nil {
	// 	panic(err)
	// }
}

func main() {
	setupDevices()
	time.Sleep(10 * time.Second)
	out, err := exec.Command("ifconfig").CombinedOutput()
	fmt.Println("[ifconfig]", string(out))
	if err != nil {
		panic(err)
	}

	interestmapsDir := "/home/pi/interestmaps/"
	interestmapsCmd := exec.Command("bash", "-c", fmt.Sprintf("cd %s; go run launch.go interestmaps-real.cfg 0 drone" , interestmapsDir))
	run(interestmapsCmd, "[interestmaps]", true, true)
	time.Sleep(20 * time.Second)

	// Turn saturatr on and off
	for {

		fmt.Printf("[saturatr] starting\n")
		saturatrCmd := exec.Command("bash", "-c", "sudo /home/pi/multisend/sender/saturatr 192.168.0.100 eth1 192.168.0.101 eth2 3.91.1.79 real")
		run(saturatrCmd, "[saturatr]", true, true)
		time.Sleep(30 * time.Second)

		fmt.Printf("[saturatr] stopping\n")
		if out, err := exec.Command("bash", "-c", "sudo killall saturatr").CombinedOutput(); err != nil {
			fmt.Printf("[saturatr] ERROR with killing process- output:%s, error:%v\n", string(out), err)
		}
		time.Sleep(1 * time.Minute)
	}
}
