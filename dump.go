package hwstats

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

// Main code from https://github.com/google/gops/blob/master/agent/agent.go

// DumpStackTrace dumps the stack trace into writer.
//   - debug = 0: gzip for go tool pprof
//   - debug >= 1: legacy text format with comments, translating addresses to function names and line numbers, so that a programmer can read the profile without tools
//   - gops used debug = 2
func DumpStackTrace(writer io.Writer, debug int) error {
	return pprof.Lookup("goroutine").WriteTo(writer, debug)
}

// DumpMutex dumps the mutex profile into writer.
//   - debug = 0: gzip for go tool pprof
//   - debug >= 1: legacy text format with comments, translating addresses to function names and line numbers, so that a programmer can read the profile without tools
//   - gops used debug = 2
func DumpMutex(writer io.Writer, debug int) error {
	return pprof.Lookup("mutex").WriteTo(writer, debug)
}

// DumpAllocs dumps the allocs profile into writer.
//   - debug = 0: gzip for go tool pprof
//   - debug >= 1: legacy text format with comments, translating addresses to function names and line numbers, so that a programmer can read the profile without tools
//   - gops used debug = 2
func DumpAllocs(writer io.Writer, debug int) error {
	runtime.SetMutexProfileFraction(1)
	defer runtime.SetMutexProfileFraction(0)
	return pprof.Lookup("allocs").WriteTo(writer, debug)
}

// DumpHeap dumps the heap profile into writer.
//   - debug = 0: gzip for go tool pprof
//   - debug >= 1: legacy text format with comments, translating addresses to function names and line numbers, so that a programmer can read the profile without tools
//   - gops used debug = 2
func DumpHeap(writer io.Writer, debug int) error {
	return pprof.Lookup("heap").WriteTo(writer, debug)
}

// DumpBlock dumps the block profile into writer.
//   - debug = 0: gzip for go tool pprof
//   - debug >= 1: legacy text format with comments, translating addresses to function names and line numbers, so that a programmer can read the profile without tools
//   - gops used debug = 2
func DumpBlock(writer io.Writer, debug int) error {
	runtime.SetBlockProfileRate(1)
	defer runtime.SetBlockProfileRate(0)
	return pprof.Lookup("block").WriteTo(writer, debug)
}

// DumpThreadCreate dumps the thread-create profile into writer.
//   - debug = 0: gzip for go tool pprof
//   - debug >= 1: legacy text format with comments, translating addresses to function names and line numbers, so that a programmer can read the profile without tools
//   - gops used debug = 2
func DumpThreadCreate(writer io.Writer, debug int) error {
	return pprof.Lookup("threadcreate").WriteTo(writer, debug)
}

// DumpHeapProfile dumps the native heap profile into writer.
func DumpHeapProfile(writer io.Writer) error {
	return pprof.WriteHeapProfile(writer)
}

// DumpCPUProfile dumps the CPU profile into writer for duration.
//   - duration = 0: 30 seconds
func DumpCPUProfile(writer io.Writer, duration time.Duration) error {
	if err := pprof.StartCPUProfile(writer); err != nil {
		return err
	}
	if duration <= 0 {
		duration = 30 * time.Second
	}
	time.Sleep(duration)
	pprof.StopCPUProfile()
	return nil
}

// DumpTraceProfile dumps the trace into writer for duration.
//   - duration = 0: 5 seconds
func DumpTraceProfile(writer io.Writer, duration time.Duration) error {
	if err := trace.Start(writer); err != nil {
		return err
	}
	if duration <= 0 {
		duration = 5 * time.Second
	}
	time.Sleep(duration)
	trace.Stop()
	return nil
}

// DumpMemory dumps the memory stats into writer.
func DumpMemory(writer io.Writer) {
	s := GetMemoryStats()

	_, _ = fmt.Fprintf(writer, "system-total-memory: %v\n", formatBytes(s.SysTotalMemory))
	_, _ = fmt.Fprintf(writer, "system-memory-usage: %v\n", formatBytes(s.SysMemoryUsage))
	_, _ = fmt.Fprintf(writer, "total-memory: %v\n", formatBytes(s.TotalMemory))
	_, _ = fmt.Fprintf(writer, "memory-usage: %v\n", formatBytes(s.MemoryUsage))
	_, _ = fmt.Fprintf(writer, "alloc: %v\n", formatBytes(s.Alloc))
	_, _ = fmt.Fprintf(writer, "total-alloc: %v\n", formatBytes(s.TotalAlloc))
	_, _ = fmt.Fprintf(writer, "sys: %v\n", formatBytes(s.Sys))
	_, _ = fmt.Fprintf(writer, "lookups: %v\n", s.Lookups)
	_, _ = fmt.Fprintf(writer, "mallocs: %v\n", s.Mallocs)
	_, _ = fmt.Fprintf(writer, "frees: %v\n", s.Frees)
	_, _ = fmt.Fprintf(writer, "heap-alloc: %v\n", formatBytes(s.HeapAlloc))
	_, _ = fmt.Fprintf(writer, "heap-sys: %v\n", formatBytes(s.HeapSys))
	_, _ = fmt.Fprintf(writer, "heap-idle: %v\n", formatBytes(s.HeapIdle))
	_, _ = fmt.Fprintf(writer, "heap-in-use: %v\n", formatBytes(s.HeapInuse))
	_, _ = fmt.Fprintf(writer, "heap-released: %v\n", formatBytes(s.HeapReleased))
	_, _ = fmt.Fprintf(writer, "heap-objects: %v\n", s.HeapObjects)
	_, _ = fmt.Fprintf(writer, "stack-in-use: %v\n", formatBytes(s.StackInuse))
	_, _ = fmt.Fprintf(writer, "stack-sys: %v\n", formatBytes(s.StackSys))
	_, _ = fmt.Fprintf(writer, "stack-mspan-inuse: %v\n", formatBytes(s.MSpanInuse))
	_, _ = fmt.Fprintf(writer, "stack-mspan-sys: %v\n", formatBytes(s.MSpanSys))
	_, _ = fmt.Fprintf(writer, "stack-mcache-inuse: %v\n", formatBytes(s.MCacheInuse))
	_, _ = fmt.Fprintf(writer, "stack-mcache-sys: %v\n", formatBytes(s.MCacheSys))
	_, _ = fmt.Fprintf(writer, "other-sys: %v\n", formatBytes(s.OtherSys))
	_, _ = fmt.Fprintf(writer, "gc-sys: %v\n", formatBytes(s.GCSys))
	_, _ = fmt.Fprintf(writer, "next-gc: when heap-alloc >= %v\n", formatBytes(s.NextGC))
	lastGC := "-"
	if s.LastGC != 0 {
		lastGC = fmt.Sprint(time.Unix(0, int64(s.LastGC)))
	}
	_, _ = fmt.Fprintf(writer, "last-gc: %v\n", lastGC)
	_, _ = fmt.Fprintf(writer, "gc-pause-total: %v\n", time.Duration(s.PauseTotalNs))
	_, _ = fmt.Fprintf(writer, "gc-pause: %v\n", s.PauseNs[(s.NumGC+255)%256])
	_, _ = fmt.Fprintf(writer, "gc-pause-end: %v\n", s.PauseEnd[(s.NumGC+255)%256])
	_, _ = fmt.Fprintf(writer, "num-gc: %v\n", s.NumGC)
	_, _ = fmt.Fprintf(writer, "num-forced-gc: %v\n", s.NumForcedGC)
	_, _ = fmt.Fprintf(writer, "gc-cpu-fraction: %v\n", s.GCCPUFraction)
	_, _ = fmt.Fprintf(writer, "enable-gc: %v\n", s.EnableGC)
	_, _ = fmt.Fprintf(writer, "debug-gc: %v\n", s.DebugGC)
}

// DumpGoroutine dumps the goroutine stats into writer.
func DumpGoroutine(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "goroutines: %v\n", runtime.NumGoroutine())
	_, _ = fmt.Fprintf(writer, "OS threads: %v\n", pprof.Lookup("threadcreate").Count())
	_, _ = fmt.Fprintf(writer, "GOMAXPROCS: %v\n", runtime.GOMAXPROCS(0))
	_, _ = fmt.Fprintf(writer, "num CPU: %v\n", runtime.NumCPU())
}

// DumpAll dumps the memory, goroutine and stack trace into files.
// write an error file named "dump-name.err" if any error occurs
func DumpAll(dir string, callback func(path string)) error {
	err := os.MkdirAll(dir, 0)
	if err != nil {
		return err
	}
	if callback == nil {
		callback = func(path string) {
			log.Printf("dump file: %v ok!", path)
		}
	}

	dumpFile(func(writer io.Writer) error {
		DumpMemory(writer)
		return nil
	}, filepath.Join(dir, "memory.txt"), callback)

	dumpFile(func(writer io.Writer) error {
		return DumpStackTrace(writer, 0)
	}, filepath.Join(dir, "stack-trace.profile"), callback)

	dumpFile(func(writer io.Writer) error {
		return DumpHeapProfile(writer)
	}, filepath.Join(dir, "heap.profile"), callback)

	dumpFile(func(writer io.Writer) error {
		DumpGoroutine(writer)
		return nil
	}, filepath.Join(dir, "goroutine.txt"), callback)

	dumpFile(func(writer io.Writer) error {
		return DumpBlock(writer, 0)
	}, filepath.Join(dir, "block.profile"), callback)

	dumpFile(func(writer io.Writer) error {
		return DumpMutex(writer, 0)
	}, filepath.Join(dir, "mutex.profile"), callback)

	dumpFile(func(writer io.Writer) error {
		return DumpThreadCreate(writer, 0)
	}, filepath.Join(dir, "thread-create.profile"), callback)

	dumpFile(func(writer io.Writer) error {
		return DumpAllocs(writer, 0)
	}, filepath.Join(dir, "allocs.profile"), callback)

	// 30s for cpu profile
	dumpFile(func(writer io.Writer) error {
		return DumpCPUProfile(writer, 0)
	}, filepath.Join(dir, "cpu-profile.profile"), callback)

	// 5s for trace profile
	dumpFile(func(writer io.Writer) error {
		return DumpTraceProfile(writer, 0)
	}, filepath.Join(dir, "trace-profile.profile"), callback)

	return nil
}

// dumpFile dumps the result of fn into path. If fn returns an error, it will be written to path.err.
func dumpFile(fn func(writer io.Writer) error, path string, callback func(path string)) {
	f, err := os.Create(path)
	if err != nil {
		_ = os.WriteFile(path+".err", []byte(err.Error()), 0)
		callback(path + ".err")
		return
	}
	defer f.Close()
	err = fn(f)
	if err != nil {
		_ = os.WriteFile(path+".err", []byte(err.Error()), 0)
		callback(path + ".err")
		return
	}
	callback(path)
}
