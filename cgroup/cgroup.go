package cgroup

// Main code from https://github.com/VictoriaMetrics/VictoriaMetrics/tree/master/lib/cgroup

// Knowledge: https://fabiokung.com/2014/03/13/memory-inside-linux-containers/

func RunInDocker() bool {
	return runInDocker()
}

// RunInCgroup returns true if the current process is in a cgroup.
// Otherwise, returns false.
func RunInCgroup() bool {
	path := CgroupPath()

	return path != "/"
}
