package cgroup

import (
	"os"
)

// CgroupPath returns the path to the cgroup of the current process.
func CgroupPath() string {
	// /proc/self/cgroup contains something like this:
	// 15:name=systemd:/
	// 14:misc:/
	// 13:rdma:/
	// 12:pids:/
	// 11:hugetlb:/
	// 10:net_prio:/
	// 9:perf_event:/
	// 8:net_cls:/
	// 7:freezer:/
	// 6:devices:/
	// 5:memory:/
	// 4:blkio:/
	// 3:cpuacct:/
	// 2:cpu:/
	// 1:cpuset:/
	// 0::/
	content, err := os.ReadFile("/proc/self/cgroup")
	if err != nil {
		return "/"
	}

	cgroupPath, err := grepFirstMatch(string(content), "0::/", 2, ":")
	if err != nil {
		return "/"
	}

	return cgroupPath
}

func runInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}
