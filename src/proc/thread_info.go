package proc

import (
	"fmt"
	"strings"
)

type Thread struct {
	pid   int
	tid   int
	name  string
	utime int
	stime int
}

func NewThread(pid, tid int) *Thread {
	t := &Thread{
		pid: pid,
		tid: tid,
	}
	err := t.refresh()
	if err != nil {
		return nil
	}
	return t
}

func (t *Thread) refresh() error {
	var (
		err        error
		threadPath string
		stat       string
	)
	threadPath = fmt.Sprintf("/proc/%v/task/%v", t.pid, t.tid)
	// read the name
	if t.name, err = cat(threadPath + "/comm"); err != nil {
		return err
	}
	t.name = strings.TrimSpace(t.name)
	// get stats
	if stat, err = cat(threadPath + "/stat"); err != nil {
		return nil
	}
	s := parseStat(stat)
	t.utime = s["utime"]
	t.stime = s["stime"]

	return nil
}

func (t *Thread) Id() int {
	return t.tid
}

func (t *Thread) Name() string {
	return t.name
}

func (t *Thread) Usage() int64 {
	return int64(t.utime + t.stime)
}

func (t *Thread) Children() map[int]StatUnit {
	return nil
}
