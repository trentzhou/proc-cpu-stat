package proc

import (
	"fmt"
	"strings"
	"time"
)

type StatUnit interface {
	Id() int
	Name() string
	Usage() int64
	Children() map[int]StatUnit
}

type Stat struct {
	parent      StatUnit
	children    map[int]StatUnit
	hz          int
	reportLines []string
	lineCount   int
	lastTime    time.Time
}

func NewStat() *Stat {
	s := Stat{}
	// get hz
	hz := getCommandOutput("getconf", "CLK_TCK")
	s.hz = atoi(strings.TrimSpace(hz))
	return &s
}

const C_ABOVE = "\033[F"
const C_CLEAR = "\033[2K"
const C_HEAD = "\r"
const C_CLEAR_SCREEN = "\033[2J"

func (s *Stat) BeginStat() {
	s.reportLines = []string{}
}

func (s *Stat) FinishStat() {
	// go above and clear
	// for i := 0; i < s.lineCount; i++ {
	// 	fmt.Print(C_ABOVE + C_CLEAR)
	// }
	fmt.Print(C_CLEAR_SCREEN)

	for _, line := range s.reportLines {
		fmt.Println(line)
	}
	s.lineCount = len(s.reportLines)
}

func (s *Stat) ReportItem(old, new StatUnit, timeDelta time.Duration) {
	seconds := float32(timeDelta) / float32(time.Second)
	rate := float32(new.Usage()-old.Usage()) * float32(100) / float32(s.hz) / seconds
	indicator := strings.Repeat("#", int(rate+0.5))
	line := fmt.Sprintf("%v %-16v %8.2f  %v", new.Id(), new.Name(), rate, indicator)
	s.reportLines = append(s.reportLines, line)
}

func (s *Stat) Update(u StatUnit) {
	var (
		now       time.Time
		timeDelta time.Duration
		children  map[int]StatUnit
	)
	now = time.Now()
	s.BeginStat()

	if s.parent != nil {
		children = u.Children()
		timeDelta = now.Sub(s.lastTime)
		s.ReportItem(s.parent, u, timeDelta)
		s.reportLines = append(s.reportLines, "-------------------------------")

		keys := getMapKeys(children)

		for _, k := range keys {
			v := children[k]
			old, ok := s.children[k]
			if ok {
				s.ReportItem(old, v, timeDelta)
			}
		}
	}
	s.FinishStat()
	s.lastTime = now
	s.parent = u
	s.children = children
}
