package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func remoteDir() string {
	return "Drone-Project/measurements/"
}

func download(folder string) {
	remotePath := fmt.Sprintf("%s/%s", remoteDir(), folder)
	downloadCmd := fmt.Sprintf("dropbox_uploader.sh download %s", remotePath)
	if out, err := exec.Command("bash", "-c", downloadCmd).CombinedOutput(); err != nil {
		fmt.Println(string(out))
		panic(err)
	}
}

func main() {
	folder := "home-10-14-1"

	download(fmt.Sprintf("combined_traces/%s", folder))
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		panic(err)
	}
	csvPath := fmt.Sprintf("%s/%s", folder, files[0].Name())
	tdCmd := fmt.Sprintf("Rscript --vanilla 3d.R %s", csvPath)
	if out, err := exec.Command("bash", "-c", tdCmd).CombinedOutput(); err != nil {
		fmt.Println(string(out))
		panic(err)
	}
}
