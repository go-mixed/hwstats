package cgroup_stats

import (
	"fmt"
	"gopkg.in/go-mixed/hwstats.v1/cgroup"
	"math"
	"testing"
)

func TestMemory(t *testing.T) {
	t.Log("TotalMemory:", prettyByteSize(TotalMemory()))
	t.Log("FreeMemory:", prettyByteSize(FreeMemory()))
}

func TestCpu(t *testing.T) {
	t.Log("AvailableCPUs:", AvailableCPUs())
}

func TestCGroup(t *testing.T) {
	t.Log("RunInDocker:", cgroup.RunInDocker())
	t.Log("RunInCgroup:", cgroup.RunInCgroup())
	if cgroup.RunInCgroup() {
		t.Logf("Cgroup path: %s", cgroup.CgroupPath())
		t.Logf("Cgroup CPUQuota: %f", cgroup.GetCPUQuota())
		t.Logf("Cgroup Memory Limit: %d", cgroup.GetMemoryLimit())
		t.Logf("Cgroup Hierarchical Memory Limit: %d", cgroup.GetHierarchicalMemoryLimit())
		memStat, err := cgroup.GetMemoryStat()
		if err != nil {
			t.Errorf("GetMemoryStat failed: %v", err)
		}
		t.Logf("Cgroup Memory Usage: %+v", memStat)
	}
}

func prettyByteSize(b uint64) string {
	bf := float64(b)
	for _, unit := range []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi"} {
		if math.Abs(bf) < 1024.0 {
			return fmt.Sprintf("%3.1f%sB", bf, unit)
		}
		bf /= 1024.0
	}
	return fmt.Sprintf("%.1fYiB", bf)
}
