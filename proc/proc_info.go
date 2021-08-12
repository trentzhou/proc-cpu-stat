package proc

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Process struct {
	pid     int
	name    string
	utime   int
	stime   int
	threads []*Thread
}

func NewProcess(pid int) *Process {
	p := Process{
		pid: pid,
	}
	err := p.refresh()
	if err != nil {
		return nil
	}
	return &p
}

func (p *Process) refresh() error {
	var (
		procPath  string
		err       error
		fileinfos []os.FileInfo
		stat      string
	)
	procPath = fmt.Sprintf("/proc/%v", p.pid)
	// read process name
	if p.name, err = cat(procPath + "/comm"); err != nil {
		return err
	}

	p.name = strings.TrimSpace(p.name)
	// get thread list
	if fileinfos, err = ioutil.ReadDir(procPath + "/task"); err != nil {
		return err
	}
	threads := []*Thread{}
	for _, f := range fileinfos {
		if f.IsDir() {
			tid := atoi(f.Name())
			thread := NewThread(p.pid, tid)
			threads = append(threads, thread)
		}
	}
	p.threads = threads

	// get stat
	if stat, err = cat(procPath + "/stat"); err != nil {
		return err
	}
	s := parseStat(stat)
	p.utime = s["utime"]
	p.stime = s["stime"]
	return nil
}

func (p *Process) Id() int {
	return p.pid
}

func (p *Process) Name() string {
	return p.name
}

func (p *Process) Usage() int64 {
	return int64(p.utime + p.stime)
}

func (p *Process) Children() map[int]StatUnit {
	r := make(map[int]StatUnit)
	for _, x := range p.threads {
		r[x.Id()] = x
	}
	return r
}
