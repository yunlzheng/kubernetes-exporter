package main

import "github.com/prometheus/client_golang/prometheus"

func addMetrics() map[string]*prometheus.GaugeVec {

	gaugeVecs := make(map[string]*prometheus.GaugeVec)

	// Node Metrics
	gaugeVecs["podStatus"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "pod_health_status",
			Help:      "HealthState of pod reported by kubernates",
		}, []string{"name", "health_state"})

	gaugeVecs["exporter"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "exporter_debug_state",
		Help:      "Debug",
	}, []string{"name", "state"})

	return gaugeVecs
}
