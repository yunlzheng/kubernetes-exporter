package main

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Exporter Sets up all the runtime and metrics
type Exporter struct {
	mutex     sync.RWMutex
	gaugeVecs map[string]*prometheus.GaugeVec
}

func newExporter() *Exporter {
	gaugeVecs := addMetrics()
	return &Exporter{
		gaugeVecs: gaugeVecs,
	}
}
