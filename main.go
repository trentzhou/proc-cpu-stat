package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/trentzhou/proc-cpu-stat/proc"
)

func main() {
	var (
		showAllThreads bool
		watch          bool
	)
	flag.BoolVar(&showAllThreads, "all", false, "Show stats for all threads")
	flag.BoolVar(&watch, "watch", false, "Watch realtime usage")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] PID\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) >= 1 {
		pidStr := flag.Args()[0]
		stat := proc.NewStat(showAllThreads, watch)
		pid, _ := strconv.Atoi(pidStr)
		for {
			p := proc.NewProcess(pid)
			if p == nil {
				fmt.Fprintf(os.Stderr, "Failed to find process %v\n", pidStr)
				os.Exit(1)
			}
			finished := stat.Update(p)
			if finished {
				break
			}
			time.Sleep(time.Second)
		}
	} else {
		flag.Usage()
	}
}
