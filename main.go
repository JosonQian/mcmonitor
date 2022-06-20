package main

import (
	"fmt"
	"os/exec"
)

var cmd *exec.Cmd

func main() {
	var out []byte
	cmd = exec.Command("echo 'stats slabs'| nc 192.168.50.114 11211")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(out))
}
