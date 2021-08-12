package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/trentzhou/proc-cpu-stat/proc"
)

func main() {
	if len(os.Args) == 2 {
		pidStr := os.Args[1]
		stat := proc.NewStat()
		pid, _ := strconv.Atoi(pidStr)
		for {
			p := proc.NewProcess(pid)
			if p == nil {
				fmt.Fprintf(os.Stderr, "Failed to find process %v\n", pidStr)
				os.Exit(1)
			}
			stat.Update(p)
			time.Sleep(time.Second)
		}
	}
}
