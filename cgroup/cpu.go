package cgroup

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetCPUQuota() float64 {
	cpuQuota, err := getCPUQuotaGeneric()
	if err != nil {
		return 0
	}
	if cpuQuota <= 0 {
		// The quota isn't set. This may be the case in multilevel containers.
		// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/685#issuecomment-674423728
		return getOnlineCPUCount()
	}
	return cpuQuota
}

// GetCPUSet returns cpuset.cpus value.
func GetCPUSet() string {
	// See https://www.kernel.org/doc/Documentation/cgroup-v1/cpusets.txt
	data, err := getFileContents("cpuset.cpus", "/sys/fs/cgroup/cpuset", "/proc/self/cgroup", "cpuset")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func getCPUQuotaGeneric() (float64, error) {
	quotaUS, err := getCPUStat("cpu.cfs_quota_us")
	if err == nil {
		periodUS, err := getCPUStat("cpu.cfs_period_us")
		if err == nil {
			return float64(quotaUS) / float64(periodUS), nil
		}
	}
	return getCPUQuotaV2("/sys/fs/cgroup", "/proc/self/cgroup")
}

func getCPUStat(statName string) (int64, error) {
	return getStatGeneric(statName, "/sys/fs/cgroup/cpu", "/proc/self/cgroup", "cpu,")
}

func getOnlineCPUCount() float64 {
	// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/685#issuecomment-674423728
	data, err := os.ReadFile("/sys/devices/system/cpu/online")
	if err != nil {
		return -1
	}
	n := float64(countCPUs(string(data)))
	if n <= 0 {
		return -1
	}
	return n
}

func getCPUQuotaV2(sysPrefix, cgroupPath string) (float64, error) {
	data, err := getFileContents("cpu.max", sysPrefix, cgroupPath, "")
	if err != nil {
		return 0, err
	}
	data = strings.TrimSpace(data)
	n, err := parseCPUMax(data)
	if err != nil {
		return 0, fmt.Errorf("cannot parse cpu.max file contents: %w", err)
	}
	return n, nil
}

// See https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v2.html#cpu
func parseCPUMax(data string) (float64, error) {
	bounds := strings.Split(data, " ")
	if len(bounds) != 2 {
		return 0, fmt.Errorf("unexpected line format: want 'quota period'; got: %s", data)
	}
	if bounds[0] == "max" {
		return -1, nil
	}
	quota, err := strconv.ParseUint(bounds[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse quota: %w", err)
	}
	period, err := strconv.ParseUint(bounds[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse period: %w", err)
	}
	return float64(quota) / float64(period), nil
}

func countCPUs(data string) int {
	data = strings.TrimSpace(data)
	n := 0
	for _, s := range strings.Split(data, ",") {
		n++
		if !strings.Contains(s, "-") {
			if _, err := strconv.Atoi(s); err != nil {
				return -1
			}
			continue
		}
		bounds := strings.Split(s, "-")
		if len(bounds) != 2 {
			return -1
		}
		start, err := strconv.Atoi(bounds[0])
		if err != nil {
			return -1
		}
		end, err := strconv.Atoi(bounds[1])
		if err != nil {
			return -1
		}
		n += end - start
	}
	return n
}
