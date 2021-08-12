package proc

import (
	"io/ioutil"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func cat(filename string) (string, error) {
	var (
		bytes []byte
		err   error
	)
	if bytes, err = ioutil.ReadFile(filename); err != nil {
		return "", err
	}
	return string(bytes), nil
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// parrseStat parses the stat from /proc/{pid}/stat
func parseStat(stat string) map[string]int {
	items := strings.Split(stat, " ")
	return map[string]int{
		"utime": atoi(items[13]),
		"stime": atoi(items[14]),
	}
}

func getCommandOutput(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func getMapKeys(m map[int]StatUnit) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// sort them
	sort.Ints(keys)
	return keys
}
