//go:build !linux && !darwin && !windows && !freebsd && !dragonfly && !netbsd && !openbsd

package hwstats

func sysTotalMemory() uint64 {
	return 0
}
func sysFreeMemory() uint64 {
	return 0
}
func sysMemoryUsage() uint64 {
	return 0
}
