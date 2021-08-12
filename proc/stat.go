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
	// display options
	showAll bool
	watch   bool

	parent      StatUnit
	children    map[int]StatUnit
	hz          int
	reportLines []string
	lineCount   int
	lastTime    time.Time
}

func NewStat(showAll, watch bool) *Stat {
	s := Stat{showAll: showAll, watch: watch}
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
	if s.watch {
		fmt.Print(C_CLEAR_SCREEN)
	}

	for _, line := range s.reportLines {
		fmt.Println(line)
	}
	s.lineCount = len(s.reportLines)
}

func (s *Stat) ReportItem(old, new StatUnit, timeDelta time.Duration) {
	seconds := float32(timeDelta) / float32(time.Second)
	usedJiffies := new.Usage() - old.Usage()
	if s.showAll || usedJiffies > 0 {
		rate := float32(usedJiffies) * float32(100) / float32(s.hz) / seconds
		indicator := strings.Repeat("#", int(rate+0.5))
		line := fmt.Sprintf("%8v %-16v %8.2f  %v", new.Id(), new.Name(), rate, indicator)
		s.reportLines = append(s.reportLines, line)
	}
}

// Update the process stat. Return true if break the loop
func (s *Stat) Update(u StatUnit) (finished bool) {
	var (
		now       time.Time
		timeDelta time.Duration
		children  map[int]StatUnit
	)
	finished = false
	now = time.Now()
	s.BeginStat()

	children = u.Children()
	if s.parent != nil {
		timeDelta = now.Sub(s.lastTime)
		s.ReportItem(s.parent, u, timeDelta)
		s.reportLines = append(s.reportLines, "-----------------------------------")

		keys := getMapKeys(children)

		for _, k := range keys {
			v := children[k]
			old, ok := s.children[k]
			if ok {
				s.ReportItem(old, v, timeDelta)
			}
		}
		s.FinishStat()

		if !s.watch {
			finished = true
		}
	}
	s.lastTime = now
	s.parent = u
	s.children = children
	return
}
