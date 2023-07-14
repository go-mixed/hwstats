//go:build freebsd || openbsd || dragonfly || netbsd

package hwstats

func sysTotalMemory() uint64 {
	s, err := sysctlUint64("hw.physmem")
	if err != nil {
		return 0
	}
	return s
}

func sysFreeMemory() uint64 {
	s, err := sysctlUint64("hw.usermem")
	if err != nil {
		return 0
	}
	return s
}

func sysMemoryUsage() uint64 {
	return sysTotalMemory() - sysFreeMemory()
}
