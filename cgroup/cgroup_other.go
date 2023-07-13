//go:build !linux

package cgroup

// CgroupPath returns the path to the cgroup of the current process.
func CgroupPath() string {
	return "/"
}

func runInDocker() bool {
	return false
}
