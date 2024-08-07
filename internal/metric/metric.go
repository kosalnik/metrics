// Package metric contains implementation of the metrics getter.
// Method GetMetrics receives the set of metrics from system.
package metric

import (
	"context"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

	"github.com/kosalnik/metrics/internal/log"
)

func GetMetrics(ctx context.Context) (map[string]float64, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	r := map[string]float64{
		"Alloc":         float64(m.Alloc),
		"BuckHashSys":   float64(m.BuckHashSys),
		"Frees":         float64(m.Frees),
		"GCCPUFraction": m.GCCPUFraction,
		"GCSys":         float64(m.GCSys),
		"HeapAlloc":     float64(m.HeapAlloc),
		"HeapIdle":      float64(m.HeapIdle),
		"HeapInuse":     float64(m.HeapInuse),
		"HeapObjects":   float64(m.HeapObjects),
		"HeapReleased":  float64(m.HeapReleased),
		"HeapSys":       float64(m.HeapSys),
		"LastGC":        float64(m.LastGC),
		"Lookups":       float64(m.Lookups),
		"MCacheInuse":   float64(m.MCacheInuse),
		"MCacheSys":     float64(m.MCacheSys),
		"MSpanInuse":    float64(m.MSpanInuse),
		"MSpanSys":      float64(m.MSpanSys),
		"Mallocs":       float64(m.Mallocs),
		"NextGC":        float64(m.NextGC),
		"NumForcedGC":   float64(m.NumForcedGC),
		"NumGC":         float64(m.NumGC),
		"OtherSys":      float64(m.OtherSys),
		"PauseTotalNs":  float64(m.PauseTotalNs),
		"StackInuse":    float64(m.StackInuse),
		"StackSys":      float64(m.StackSys),
		"Sys":           float64(m.Sys),
		"TotalAlloc":    float64(m.TotalAlloc),
	}

	if cpuUsage, err := cpu.PercentWithContext(ctx, 0, false); err == nil {
		r["CPUutilization1"] = float64(cpuUsage[0])
	} else {
		log.Error().Err(err).Msg("get cpu usage fail")
	}

	if memUsage, err := mem.VirtualMemory(); err == nil {
		r["TotalMemory"] = float64(memUsage.Total)
		r["FreeMemory"] = float64(memUsage.Free)
	} else {
		log.Error().Err(err).Msg("get memory usage fail")
	}

	return r, nil
}
