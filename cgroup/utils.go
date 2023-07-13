package cgroup

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

// get stats from cgroup
//   - statName: /sys/fs/cgroup/cgroup-subpath/statFileName
//   - sysfsPrefix: path to /sys/fs/cgroup/
//   - cgroupPath: path to /proc/self/cgroup
//   - cgroupGrepLine: line to grep from cgroupPath
func getStatGeneric(statFileName, sysfsPrefix, cgroupPath, cgroupGrepLine string) (int64, error) {
	data, err := getFileContents(statFileName, sysfsPrefix, cgroupPath, cgroupGrepLine)
	if err != nil {
		return 0, err
	}
	data = strings.TrimSpace(data)
	n, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse %q: %w", cgroupPath, err)
	}
	return n, nil
}

func getFileContents(statFileName, sysfsPrefix, cgroupPath, cgroupGrepLine string) (string, error) {
	// try to read "/sys/fs/cgroup/statFileName" first, eg: in docker.
	// then try to read "/sys/fs/cgroup/cgroup-subpath/statFileName", eg: in host os.
	filepath := path.Join(sysfsPrefix, statFileName)
	data, err := os.ReadFile(filepath)
	if err == nil {
		return string(data), nil
	}

	// parse cgroup-subpath from cgroupPath
	cgroupData, err := os.ReadFile(cgroupPath)
	if err != nil {
		return "", err
	}
	subPath, err := grepFirstMatch(string(cgroupData), cgroupGrepLine, 2, ":")
	if err != nil {
		return "", fmt.Errorf("cannot find cgroup path for %q in %q: %w", cgroupGrepLine, cgroupPath, err)
	}
	// sys/fs/cgroup/cgroup-subpath/statFileName
	filepath = path.Join(sysfsPrefix, subPath, statFileName)

	data, err = os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// grepFirstMatch searches match line at data and returns item from it by index with given delimiter.
func grepFirstMatch(data string, match string, index int, delimiter string) (string, error) {
	lines := strings.Split(string(data), "\n")
	for _, s := range lines {
		if !strings.Contains(s, match) {
			continue
		}
		parts := strings.Split(s, delimiter)
		if index < len(parts) {
			return strings.TrimSpace(parts[index]), nil
		}
	}
	return "", fmt.Errorf("cannot find %q in %q", match, data)
}
