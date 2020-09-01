package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"time"
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

func main() {
	monitorDir := "/home/ubuntu/data/drone-monitor/"
	fmt.Println("Starting drone monitor")
	monitorCmd := exec.Command(monitorDir + "drone-monitor")
	run(monitorCmd, "[monitor]", true, true)
	fmt.Println("Starting interestmaps")
	time.Sleep(2 * time.Second)

	frontendDir := "/home/ubuntu/data/drone-monitor/client/simple/build"
	frontendCmd := exec.Command("serve", frontendDir)
	run(frontendCmd, "[frontend]", true, true)
	time.Sleep(2 * time.Second)

	interestmapsDir := "/home/ubuntu/interestmaps/"
	interestmapsCmd := exec.Command(interestmapsDir+"interestmaps", interestmapsDir+"interestmaps-real.cfg")
	run(interestmapsCmd, "[interestmaps]", true, true)
	fmt.Println("Starting saturatr")
	time.Sleep(2 * time.Second)

	saturatrDir := "/home/ubuntu/data/multisend/sender/"
	saturatrCmd := exec.Command(saturatrDir+"saturatr", "real")
	run(saturatrCmd, "[saturatr]", true, true)

	monitorCmd.Wait()
	frontendCmd.Wait()
	interestmapsCmd.Wait()
	saturatrCmd.Wait()
}
