package hwstats

import (
	"gopkg.in/go-mixed/hwstats.v1/cgroup"
	"runtime"
)

// Main code from https://github.com/pbnjay/memory

// SysTotalMemory returns the total accessible system memory in bytes.
//
// The total accessible memory is installed physical memory size minus reserved
// areas for the kernel and hardware, if such reservations are reported by
// the operating system.
//
// If accessible memory size could not be determined, then 0 is returned.
func SysTotalMemory() uint64 {
	return sysTotalMemory()
}

// SysFreeMemory returns the total free system memory in bytes.
//
// The total free memory is installed physical memory size minus reserved
// areas for other applications running on the same system.
//
// If free memory size could not be determined, then 0 is returned.
func SysFreeMemory() uint64 {
	return sysFreeMemory()
}

// SysMemoryUsage returns the total used system memory in bytes.
//
// The total used memory is installed physical memory size minus reserved
// areas for other applications running on the same system.
//
// If used memory size could not be determined, then 0 is returned.
func SysMemoryUsage() uint64 {
	return sysMemoryUsage()
}

// TotalMemory returns the really total memory, if run in cgroup, it will return
// the cgroup memory limit, otherwise it will return the system total memory
func TotalMemory() uint64 {
	totalMemory := SysTotalMemory()

	if cgroup.RunInCgroup() {
		if cgroupMemoryLimit := uint64(cgroup.GetMemoryLimit()); cgroupMemoryLimit > totalMemory {
			return totalMemory
		} else {
			return cgroupMemoryLimit
		}
	}

	return totalMemory
}

// MemoryUsage returns the real memory usage, if run in cgroup, it will return
// the cgroup memory RSS+Cache usage, otherwise it will return the system memory usage
func MemoryUsage() uint64 {
	if cgroup.RunInCgroup() {
		if memStat, err := cgroup.GetMemoryStat(); err != nil {
			return 0
		} else {
			return uint64(memStat.Rss + memStat.Cache)
		}
	} else {
		return SysMemoryUsage()
	}
}

type MemoryStats struct {
	// MemStats is the memory statistics of current process.
	runtime.MemStats
	// SysTotalMemory is the total accessible system memory in bytes.
	SysTotalMemory uint64 `json:"sys_total_memory" yaml:"sys_total_memory"`
	// SysMemoryUsage is the total used system memory in bytes.
	SysMemoryUsage uint64 `json:"sys_memory_usage" yaml:"sys_memory_usage"`
	// TotalMemory is the really total memory, cgroup memory limit or system total memory.
	TotalMemory uint64 `json:"total_memory" yaml:"total_memory"`
	// MemoryUsage is the real memory usage, cgroup memory usage or system memory usage.
	MemoryUsage uint64 `json:"memory_usage" yaml:"memory_usage"`
}

// GetMemoryStats returns the memory statistics of system,and the current process.
func GetMemoryStats() MemoryStats {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	return MemoryStats{
		MemStats:       ms,
		SysTotalMemory: SysFreeMemory(),
		SysMemoryUsage: SysMemoryUsage(),
		TotalMemory:    TotalMemory(),
		MemoryUsage:    MemoryUsage(),
	}
}
