package utils

import (
	"bufio"
	"os/exec"
)

func RunCmd(cmd *exec.Cmd, tag string, onStdout func(string), onStderr func(string)) {
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	go func() {
		buf := bufio.NewReader(stdout)
		for {
			line, _, err := buf.ReadLine()
			if err != nil {
				break
			}
			onStdout(string(line))
		}
	}()

	go func() {
		buf := bufio.NewReader(stderr)
		for {
			line, _, err := buf.ReadLine()
			if err != nil {
				break
			}
			onStderr(string(line))
		}
	}()
}
