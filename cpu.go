package cgroup_stats

import (
	"os"
	"runtime"
)

// AvailableCPUs returns the number of available CPU cores for the app.
//
// The number is rounded to the next integer value if fractional number of CPU cores are available.
func AvailableCPUs() int {
	return runtime.GOMAXPROCS(-1)
}

// UpdateGOMAXPROCSToCPUQuota updates GOMAXPROCS to cpuQuota if GOMAXPROCS isn't set in environment var.
func UpdateGOMAXPROCSToCPUQuota(cpuQuota float64) {
	if v := os.Getenv("GOMAXPROCS"); v != "" {
		// Do not override explicitly set GOMAXPROCS.
		return
	}
	gomaxprocs := int(cpuQuota + 0.5)
	numCPU := runtime.NumCPU()
	if gomaxprocs > numCPU {
		// There is no sense in setting more GOMAXPROCS than the number of available CPU cores.
		gomaxprocs = numCPU
	}
	if gomaxprocs <= 0 {
		gomaxprocs = 1
	}
	runtime.GOMAXPROCS(gomaxprocs)
}
