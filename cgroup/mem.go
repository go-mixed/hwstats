package cgroup

import (
	"strconv"
	"strings"
)

// GetMemoryLimit returns cgroup memory limit
func GetMemoryLimit() int64 {
	// Try determining the amount of memory inside docker container.
	// See https://stackoverflow.com/questions/42187085/check-mem-limit-within-a-docker-container
	//
	// Read memory limit according to https://unix.stackexchange.com/questions/242718/how-to-find-out-how-much-memory-lxc-container-is-allowed-to-consume
	// This should properly determine the limit inside lxc container.
	// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/84
	n, err := getMemStat("memory.limit_in_bytes")
	if err == nil {
		return n
	}
	n, err = getMemStatV2("memory.max")
	if err != nil {
		return 0
	}
	return n
}

func getMemStatV2(statName string) (int64, error) {
	// See https: //www.kernel.org/doc/html/latest/admin-guide/cgroup-v2.html#memory-interface-files
	return getStatGeneric(statName, "/sys/fs/cgroup", "/proc/self/cgroup", "")
}

func getMemStat(statName string) (int64, error) {
	return getStatGeneric(statName, "/sys/fs/cgroup/memory", "/proc/self/cgroup", "memory")
}

// GetHierarchicalMemoryLimit returns hierarchical memory limit
// https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
func GetHierarchicalMemoryLimit() int64 {
	// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/699
	n, err := getHierarchicalMemoryLimit("/sys/fs/cgroup/memory", "/proc/self/cgroup")
	if err != nil {
		return 0
	}
	return n
}

func getHierarchicalMemoryLimit(sysfsPrefix, cgroupPath string) (int64, error) {
	data, err := getFileContents("memory.stat", sysfsPrefix, cgroupPath, "memory")
	if err != nil {
		return 0, err
	}
	memStat, err := grepFirstMatch(data, "hierarchical_memory_limit", 1, " ")
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(memStat, 10, 64)
}

type MemoryStat struct {
	// See https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt for description.
	Cache                   int64 `json:"cache" yaml:"cache" mapstructure:"cache"`
	Rss                     int64 `json:"rss" yaml:"rss" mapstructure:"rss"`
	RssHuge                 int64 `json:"rss_huge" yaml:"rss_huge" mapstructure:"rss_huge"`
	Shmem                   int64 `json:"shmem" yaml:"shmem" mapstructure:"shmem"`
	MappedFile              int64 `json:"mapped_file" yaml:"mapped_file" mapstructure:"mapped_file"`
	Dirty                   int64 `json:"dirty" yaml:"dirty" mapstructure:"dirty"`
	Writeback               int64 `json:"writeback" yaml:"writeback" mapstructure:"writeback"`
	Swap                    int64 `json:"swap" yaml:"swap" mapstructure:"swap"`
	Pgpgin                  int64 `json:"pgpgin" yaml:"pgpgin" mapstructure:"pgpgin"`
	Pgpgout                 int64 `json:"pgpgout" yaml:"pgpgout" mapstructure:"pgpgout"`
	Pgfault                 int64 `json:"pgfault" yaml:"pgfault" mapstructure:"pgfault"`
	Pgmajfault              int64 `json:"pgmajfault" yaml:"pgmajfault" mapstructure:"pgmajfault"`
	InactiveAnon            int64 `json:"inactive_anon" yaml:"inactive_anon" mapstructure:"inactive_anon"`
	ActiveAnon              int64 `json:"active_anon" yaml:"active_anon" mapstructure:"active_anon"`
	InactiveFile            int64 `json:"inactive_file" yaml:"inactive_file" mapstructure:"inactive_file"`
	ActiveFile              int64 `json:"active_file" yaml:"active_file" mapstructure:"active_file"`
	Unevictable             int64 `json:"unevictable" yaml:"unevictable" mapstructure:"unevictable"`
	HierarchicalMemoryLimit int64 `json:"hierarchical_memory_limit" yaml:"hierarchical_memory_limit" mapstructure:"hierarchical_memory_limit"`
	HierarchicalMemswLimit  int64 `json:"hierarchical_memsw_limit" yaml:"hierarchical_memsw_limit" mapstructure:"hierarchical_memsw_limit"`
	TotalCache              int64 `json:"total_cache" yaml:"total_cache" mapstructure:"total_cache"`
	TotalRss                int64 `json:"total_rss" yaml:"total_rss" mapstructure:"total_rss"`
	TotalRssHuge            int64 `json:"total_rss_huge" yaml:"total_rss_huge" mapstructure:"total_rss_huge"`
	TotalShmem              int64 `json:"total_shmem" yaml:"total_shmem" mapstructure:"total_shmem"`
	TotalMappedFile         int64 `json:"total_mapped_file" yaml:"total_mapped_file" mapstructure:"total_mapped_file"`
	TotalDirty              int64 `json:"total_dirty" yaml:"total_dirty" mapstructure:"total_dirty"`
	TotalWriteback          int64 `json:"total_writeback" yaml:"total_writeback" mapstructure:"total_writeback"`
	TotalSwap               int64 `json:"total_swap" yaml:"total_swap" mapstructure:"total_swap"`
	TotalPgpgin             int64 `json:"total_pgpgin" yaml:"total_pgpgin" mapstructure:"total_pgpgin"`
	TotalPgpgout            int64 `json:"total_pgpgout" yaml:"total_pgpgout" mapstructure:"total_pgpgout"`
	TotalPgfault            int64 `json:"total_pgfault" yaml:"total_pgfault" mapstructure:"total_pgfault"`
	TotalPgmajfault         int64 `json:"total_pgmajfault" yaml:"total_pgmajfault" mapstructure:"total_pgmajfault"`
	TotalInactiveAnon       int64 `json:"total_inactive_anon" yaml:"total_inactive_anon" mapstructure:"total_inactive_anon"`
	TotalActiveAnon         int64 `json:"total_active_anon" yaml:"total_active_anon" mapstructure:"total_active_anon"`
	TotalInactiveFile       int64 `json:"total_inactive_file" yaml:"total_inactive_file" mapstructure:"total_inactive_file"`
	TotalActiveFile         int64 `json:"total_active_file" yaml:"total_active_file" mapstructure:"total_active_file"`
	TotalUnevictable        int64 `json:"total_unevictable" yaml:"total_unevictable" mapstructure:"total_unevictable"`
}

// GetMemoryStat returns memory statistics for the current process.
func GetMemoryStat() (*MemoryStat, error) {
	data, err := getFileContents("memory.stat", "/sys/fs/cgroup/memory", "/proc/self/cgroup", "memory")
	if err != nil {
		return nil, err
	}

	m := map[string]int64{}
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		segment := strings.Split(line, " ")
		m[segment[0]], _ = strconv.ParseInt(segment[1], 10, 64)
	}

	return &MemoryStat{
		Cache:                   m["cache"],
		Rss:                     m["rss"],
		RssHuge:                 m["rss_huge"],
		Shmem:                   m["shmem"],
		MappedFile:              m["mapped_file"],
		Dirty:                   m["dirty"],
		Writeback:               m["writeback"],
		Swap:                    m["swap"],
		Pgpgin:                  m["pgpgin"],
		Pgpgout:                 m["pgpgout"],
		Pgfault:                 m["pgfault"],
		Pgmajfault:              m["pgmajfault"],
		InactiveAnon:            m["inactive_anon"],
		ActiveAnon:              m["active_anon"],
		InactiveFile:            m["inactive_file"],
		ActiveFile:              m["active_file"],
		Unevictable:             m["unevictable"],
		HierarchicalMemoryLimit: m["hierarchical_memory_limit"],
		HierarchicalMemswLimit:  m["hierarchical_memsw_limit"],
		TotalCache:              m["total_cache"],
		TotalRss:                m["total_rss"],
		TotalRssHuge:            m["total_rss_huge"],
		TotalShmem:              m["total_shmem"],
		TotalMappedFile:         m["total_mapped_file"],
		TotalDirty:              m["total_dirty"],
		TotalWriteback:          m["total_writeback"],
		TotalSwap:               m["total_swap"],
		TotalPgpgin:             m["total_pgpgin"],
		TotalPgpgout:            m["total_pgpgout"],
		TotalPgfault:            m["total_pgfault"],
		TotalPgmajfault:         m["total_pgmajfault"],
		TotalInactiveAnon:       m["total_inactive_anon"],
		TotalActiveAnon:         m["total_active_anon"],
		TotalInactiveFile:       m["total_inactive_file"],
		TotalActiveFile:         m["total_active_file"],
		TotalUnevictable:        m["total_unevictable"],
	}, nil
}
