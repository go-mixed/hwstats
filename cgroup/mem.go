package cgroup

import (
	"strconv"
	"strings"
)

// GetMemoryLimit returns cgroup memory limit from "memory.limit_in_bytes" file.
func GetMemoryLimit() int64 {
	// Try determining the amount of memory inside docker container.
	// See https://stackoverflow.com/questions/42187085/check-mem-limit-within-a-docker-container
	//
	// Read memory limit according to https://unix.stackexchange.com/questions/242718/how-to-find-out-how-much-memory-lxc-container-is-allowed-to-consume
	// This should properly determine the limit inside lxc container.
	// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/84
	return getMemStat("memory.limit_in_bytes", "memory.max")
}

// GetMemoryUsage returns memory usage from "memory.usage_in_bytes" file.
func GetMemoryUsage() int64 {
	return getMemStat("memory.usage_in_bytes", "memory.current")
}

// GetMemoryFailcnt returns memory failcnt from "memory.failcnt" file.
func GetMemoryFailcnt() int64 {
	return getMemStat("memory.failcnt", "")
}

// GetMemoryMaxUsage returns maximum memory usage from "memory.max_usage_in_bytes" file.
func GetMemoryMaxUsage() int64 {
	return getMemStat("memory.max_usage_in_bytes", "memory.max_usage")
}

// GetMemoryHierarchicalLimit returns hierarchical memory limit from "memory.hierarchical_memory_limit" file.
func GetMemoryHierarchicalLimit() int64 {
	return getMemStat("memory.hierarchical_memory_limit", "memory.high")
}

func GetMemoryOOMControl() int64 {
	return getMemStat("memory.oom_control", "memory.oom_kill_disable")
}

// GetHierarchicalMemoryLimit returns hierarchical memory limit from "memory.stat" file.
// https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
func GetHierarchicalMemoryLimit() int64 {
	// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/699
	data, err := getFileContents("memory.stat", "/sys/fs/cgroup/memory", "/proc/self/cgroup", "memory")
	if err != nil {
		return 0
	}
	memStat, err := grepFirstMatch(data, "hierarchical_memory_limit", 1, " ")
	if err != nil {
		return 0
	}
	n, _ := strconv.ParseInt(memStat, 10, 64)
	return n
}

func getMemStat(statFileName string, v2StatFileName string) int64 {
	n, err := getStatGeneric(statFileName, "/sys/fs/cgroup/memory", "/proc/self/cgroup", "memory")
	if err == nil {
		return n
	}

	if v2StatFileName == "" {
		return 0
	}

	// See https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v2.html#memory-interface-files
	n, err = getStatGeneric(v2StatFileName, "/sys/fs/cgroup", "/proc/self/cgroup", "")
	if err != nil {
		return 0
	}
	return n
}

// MemoryStat https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
type MemoryStat struct {
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
		segments := strings.Split(line, " ")
		if len(segments) >= 2 {
			m[segments[0]], _ = strconv.ParseInt(segments[1], 10, 64)
		}
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

// MemoryStatV2 https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v2.html#memory.stat
type MemoryStatV2 struct {
	Anon                   int64 `json:"anon" yaml:"anon" mapstructure:"anon"`
	File                   int64 `json:"file" yaml:"file" mapstructure:"file"`
	Kernel                 int64 `json:"kernel" yaml:"kernel" mapstructure:"kernel"`
	KernelStack            int64 `json:"kernel_stack" yaml:"kernel_stack" mapstructure:"kernel_stack"`
	PageTables             int64 `json:"pagetables" yaml:"pagetables" mapstructure:"pagetables"`
	SecPageTables          int64 `json:"sec_pagetables" yaml:"sec_pagetables" mapstructure:"sec_pagetables"`
	PerCPU                 int64 `json:"percpu" yaml:"per_cpu" mapstructure:"percpu"`
	Sock                   int64 `json:"sock" yaml:"sock" mapstructure:"sock"`
	Shmem                  int64 `json:"shmem" yaml:"shmem" mapstructure:"shmem"`
	ZSwap                  int64 `json:"zswap" yaml:"zswap" mapstructure:"zswap"`
	ZSwapped               int64 `json:"zswapped" yaml:"zswapped" mapstructure:"zswapped"`
	FileMapped             int64 `json:"file_mapped" yaml:"file_mapped" mapstructure:"file_mapped"`
	FileDirty              int64 `json:"file_dirty" yaml:"file_dirty" mapstructure:"file_dirty"`
	FileWriteback          int64 `json:"file_writeback" yaml:"file_writeback" mapstructure:"file_writeback"`
	SwapCached             int64 `json:"swapcached" yaml:"swapcached" mapstructure:"swapcached"`
	AnonThp                int64 `json:"anon_thp" yaml:"anon_thp" mapstructure:"anon_thp"`
	FileThp                int64 `json:"file_thp" yaml:"file_thp" mapstructure:"file_thp"`
	ShmemThp               int64 `json:"shmem_thp" yaml:"shmem_thp" mapstructure:"shmem_thp"`
	InactiveAnon           int64 `json:"inactive_anon" yaml:"inactive_anon" mapstructure:"inactive_anon"`
	ActiveAnon             int64 `json:"active_anon" yaml:"active_anon" mapstructure:"active_anon"`
	InactiveFile           int64 `json:"inactive_file" yaml:"inactive_file" mapstructure:"inactive_file"`
	ActiveFile             int64 `json:"active_file" yaml:"active_file" mapstructure:"active_file"`
	SlabReclaimable        int64 `json:"slab_reclaimable" yaml:"slab_reclaimable" mapstructure:"slab_reclaimable"`
	SlabUnreclaimable      int64 `json:"slab_unreclaimable" yaml:"slab_unreclaimable" mapstructure:"slab_unreclaimable"`
	Slab                   int64 `json:"slab" yaml:"slab" mapstructure:"slab"`
	WorkingsetRefaultAnon  int64 `json:"workingset_refault_anon" yaml:"workingset_refault_anon" mapstructure:"workingset_refault_anon"`
	WorkingsetRefaultFile  int64 `json:"workingset_refault_file" yaml:"workingset_refault_file" mapstructure:"workingset_refault_file"`
	WorkingsetActivateAnon int64 `json:"workingset_activate_anon" yaml:"workingset_activate_anon" mapstructure:"workingset_activate_anon"`
	WorkingsetActivateFile int64 `json:"workingset_activate_file" yaml:"workingset_activate_file" mapstructure:"workingset_activate_file"`
	WorkingsetRestoreAnon  int64 `json:"workingset_restore_anon" yaml:"workingset_restore_anon" mapstructure:"workingset_restore_anon"`
	WorkingsetRestoreFile  int64 `json:"workingset_restore_file" yaml:"workingset_restore_file" mapstructure:"workingset_restore_file"`
	WorkingsetNodereclaim  int64 `json:"workingset_nodereclaim" yaml:"workingset_nodereclaim" mapstructure:"workingset_nodereclaim"`
	PgScan                 int64 `json:"pgscan" yaml:"pgscan" mapstructure:"pgscan"`
	PgSteal                int64 `json:"pgsteal" yaml:"pgsteal" mapstructure:"pgsteal"`
	PgScanKswapd           int64 `json:"pgscan_kswapd" yaml:"pgscan_kswapd" mapstructure:"pgscan_kswapd"`
	PgScanDirect           int64 `json:"pgscan_direct" yaml:"pgscan_direct" mapstructure:"pgscan_direct"`
	PgScanKhugepaged       int64 `json:"pgscan_khugepaged" yaml:"pgscan_khugepaged" mapstructure:"pgscan_khugepaged"`
	PgStealKswapd          int64 `json:"pgsteal_kswapd" yaml:"pgsteal_kswapd" mapstructure:"pgsteal_kswapd"`
	PgStealDirect          int64 `json:"pgsteal_direct" yaml:"pgsteal_direct" mapstructure:"pgsteal_direct"`
	PgStealKhugepaged      int64 `json:"pgsteal_khugepaged" yaml:"pgsteal_khugepaged" mapstructure:"pgsteal_khugepaged"`
	PgFault                int64 `json:"pgfault" yaml:"pgfault" mapstructure:"pgfault"`
	PgMajFault             int64 `json:"pgmajfault" yaml:"pgmajfault" mapstructure:"pgmajfault"`
	PgRefill               int64 `json:"pgrefill" yaml:"pgrefill" mapstructure:"pgrefill"`
	PgActivate             int64 `json:"pgactivate" yaml:"pgactivate" mapstructure:"pgactivate"`
	PgDeactivate           int64 `json:"pgdeactivate" yaml:"pgdeactivate" mapstructure:"pgdeactivate"`
	PgLazyfree             int64 `json:"pglazyfree" yaml:"pglazyfree" mapstructure:"pglazyfree"`
	PgLazyfreed            int64 `json:"pglazyfreed" yaml:"pglazyfreed" mapstructure:"pglazyfreed"`
	ThpFaultAlloc          int64 `json:"thp_fault_alloc" yaml:"thp_fault_alloc" mapstructure:"thp_fault_alloc"`
	ThpCollapseAlloc       int64 `json:"thp_collapse_alloc" yaml:"thp_collapse_alloc" mapstructure:"thp_collapse_alloc"`
}

func GetMemoryStatV2() (*MemoryStatV2, error) {
	data, err := getFileContents("memory.stat", "/sys/fs/cgroup", "/proc/self/cgroup", "")
	if err != nil {
		return nil, err
	}

	m := map[string]int64{}
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		segments := strings.Split(line, " ")
		if len(segments) >= 2 {
			m[segments[0]], _ = strconv.ParseInt(segments[1], 10, 64)
		}
	}

	return &MemoryStatV2{
		Anon:                   m["anon"],
		File:                   m["file"],
		Kernel:                 m["kernel"],
		KernelStack:            m["kernel_stack"],
		PageTables:             m["pagetables"],
		SecPageTables:          m["secpagetables"],
		PerCPU:                 m["percpu"],
		Sock:                   m["sock"],
		Shmem:                  m["shmem"],
		ZSwap:                  m["zswap"],
		ZSwapped:               m["zswapped"],
		FileMapped:             m["file_mapped"],
		FileDirty:              m["file_dirty"],
		FileWriteback:          m["file_writeback"],
		SwapCached:             m["swapcached"],
		AnonThp:                m["anon_thp"],
		FileThp:                m["file_thp"],
		ShmemThp:               m["shmem_thp"],
		InactiveAnon:           m["inactive_anon"],
		ActiveAnon:             m["active_anon"],
		InactiveFile:           m["inactive_file"],
		ActiveFile:             m["active_file"],
		SlabReclaimable:        m["slab_reclaimable"],
		SlabUnreclaimable:      m["slab_unreclaimable"],
		Slab:                   m["slab"],
		WorkingsetRefaultAnon:  m["workingset_refault_anon"],
		WorkingsetRefaultFile:  m["workingset_refault_file"],
		WorkingsetActivateAnon: m["workingset_activate_anon"],
		WorkingsetActivateFile: m["workingset_activate_file"],
		WorkingsetRestoreAnon:  m["workingset_restore_anon"],
		WorkingsetRestoreFile:  m["workingset_restore_file"],
		WorkingsetNodereclaim:  m["workingset_nodereclaim"],
		PgScan:                 m["pgscan"],
		PgSteal:                m["pgsteal"],
		PgScanKswapd:           m["pgscan_kswapd"],
		PgScanDirect:           m["pgscan_direct"],
		PgScanKhugepaged:       m["pgscan_khugepaged"],
		PgStealKswapd:          m["pgsteal_kswapd"],
		PgStealDirect:          m["pgsteal_direct"],
		PgStealKhugepaged:      m["pgsteal_khugepaged"],
		PgFault:                m["pgfault"],
		PgMajFault:             m["pgmajfault"],
		PgRefill:               m["pgrefill"],
		PgActivate:             m["pgactivate"],
		PgDeactivate:           m["pgdeactivate"],
		PgLazyfree:             m["pglazyfree"],
		PgLazyfreed:            m["pglazyfreed"],
		ThpFaultAlloc:          m["thp_fault_alloc"],
		ThpCollapseAlloc:       m["thp_collapse_alloc"],
	}, nil
}
