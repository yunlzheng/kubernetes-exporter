package main

import "github.com/prometheus/client_golang/prometheus"

func addMetrics() map[string]*prometheus.GaugeVec {

	gaugeVecs := make(map[string]*prometheus.GaugeVec)

	// Node Metrics
	gaugeVecs["nodes"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "node_status",
			Help:      "status of node reported by kubernates",
		}, []string{"name", "namespace"})

	gaugeVecs["deployments"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "deployment_status",
			Help:      "status of deployment reported by kubernates",
		}, []string{"name", "namespace"})

	gaugeVecs["stacks"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "stack_status",
			Help:      "status of stack reported by kubernates",
		}, []string{"name", "namespace"})

	gaugeVecs["pods"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "pod_status",
			Help:      "status of pod reported by kubernates",
		}, []string{"name", "namespace", "podPhase", "hostIP", "podIP", "reason", "message", "containerSize"})

	return gaugeVecs
}
