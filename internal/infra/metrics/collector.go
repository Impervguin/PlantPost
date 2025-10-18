package metrics

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	clkTck             float64
	cpuReadings        []float64
	memoryReadings     []uint64
	cpuCurrentValue    float64
	cpuMaxValue        float64
	cpuMinValue        float64
	cpuAvgValue        float64
	memoryCurrentValue uint64
	memoryMaxValue     uint64
	memoryMinValue     uint64
	memoryAvgValue     uint64
	lastUse            float64
	lastTotal          float64

	once sync.Once
)

func getMemoryValue() uint64 {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem.HeapAlloc + mem.StackInuse
}

func getCPUSample() (use, total float64) {
	pid := os.Getpid()
	contents, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return
	}
	infos := strings.Split(string(contents), " ")
	utime, stime := infos[13], infos[14]

	contents, err = os.ReadFile("/proc/uptime")
	if err != nil {
		return
	}
	uptime := strings.Split(string(contents), " ")[0]

	uptimeFloat, err := strconv.ParseFloat(uptime, 64)
	if err != nil {
		return 0, 0
	}
	uTimeTicks, err := strconv.ParseFloat(utime, 64)
	if err != nil {
		return 0, 0
	}
	sTimeTicks, err := strconv.ParseFloat(stime, 64)
	if err != nil {
		return 0, 0
	}

	useTime := (uTimeTicks + sTimeTicks) / clkTck

	return useTime, uptimeFloat
}

func getCpuValue() float64 {
	use, total := getCPUSample()
	if lastUse == 0 || lastTotal == 0 {
		lastUse = use
		lastTotal = total
		return 0
	}
	useDiff := use - lastUse
	totalDiff := total - lastTotal
	percent := (useDiff / totalDiff) * 100
	if percent > 100 {
		percent = 100
	}
	lastUse = use
	lastTotal = total
	return percent
}

func updateCPUMetrics() float64 {
	percent := getCpuValue()
	cpuCurrentValue = percent
	if percent > cpuMaxValue {
		cpuMaxValue = percent
	}
	if percent < cpuMinValue {
		cpuMinValue = percent
	}
	cpuReadings = append(cpuReadings, percent)
	cpuAvgValue = (cpuAvgValue*float64(len(cpuReadings)-1) + percent) / float64(len(cpuReadings))
	return percent
}

func updateMemoryMetrics() uint64 {
	mem := getMemoryValue()
	memoryCurrentValue = mem
	if mem > memoryMaxValue {
		memoryMaxValue = mem
	}
	if mem < memoryMinValue {
		memoryMinValue = mem
	}
	memoryReadings = append(memoryReadings, mem)
	memoryAvgValue = (memoryAvgValue*uint64(len(memoryReadings)-1) + mem) / uint64(len(memoryReadings))
	return mem
}

func setClockTicks() error {
	res, err := exec.Command("getconf", "CLK_TCK").Output()
	if err != nil {
		return err
	}
	clkTck, err = strconv.ParseFloat(strings.Trim(string(res), "\n"), 64)
	if err != nil {
		return err
	}
	return nil
}

func RegisterMetricCollector(sleepMs int) error {
	once.Do(func() {
		err := setClockTicks()
		if err != nil {
			panic(err)
		}
		cpuReadings = make([]float64, 0)
		memoryReadings = make([]uint64, 0)
		go func() {
			for {
				updateCPUMetrics()
				updateMemoryMetrics()
				time.Sleep(time.Duration(sleepMs) * time.Millisecond)
			}
		}()
	})
	return nil
}

type MetricsOutput struct {
	CpuCurrentValue    float64 `json:"cpu_current_value"`
	CpuMaxValue        float64 `json:"cpu_max_value"`
	CpuMinValue        float64 `json:"cpu_min_value"`
	CpuAvgValue        float64 `json:"cpu_avg_value"`
	MemoryCurrentValue uint64  `json:"memory_current_value"`
	MemoryMaxValue     uint64  `json:"memory_max_value"`
	MemoryMinValue     uint64  `json:"memory_min_value"`
	MemoryAvgValue     uint64  `json:"memory_avg_value"`
}

func GetMetrics() MetricsOutput {
	once.Do(func() {
		panic("metrics not initialized")
	})
	return MetricsOutput{
		CpuCurrentValue:    cpuCurrentValue,
		CpuMaxValue:        cpuMaxValue,
		CpuMinValue:        cpuMinValue,
		CpuAvgValue:        cpuAvgValue,
		MemoryCurrentValue: memoryCurrentValue,
		MemoryMaxValue:     memoryMaxValue,
		MemoryMinValue:     memoryMinValue,
		MemoryAvgValue:     memoryAvgValue,
	}
}
