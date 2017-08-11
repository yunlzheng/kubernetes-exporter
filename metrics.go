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
		}, []string{
			"name",
			"nodeState",
			"osImage",
			"containerRuntimeVersion",
			"kubeletVersion",
			"operatingSystem",
			"architecture",
			"hostname",
			"externalIp",
			"internalIp",
		})

	gaugeVecs["components"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "component_status",
			Help:      "status of kubenetes component reported by kubernates",
		}, []string{"name", "namespace"})

	gaugeVecs["stacks"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "stacks_status",
			Help:      "status of stacks reported by kubernates",
		}, []string{"name", "namespace"})

	gaugeVecs["deployments"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "deployment_status",
			Help:      "status of deployment reported by kubernates",
		}, []string{"name", "namespace"})

	gaugeVecs["statefulsets"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "stateful_sets_status",
			Help:      "status of stateful sets reported by kubernates",
		}, []string{"name", "namespace"})

	gaugeVecs["daemonsets"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "daemon_sets_status",
			Help:      "status of daemon sets reported by kubernates",
		}, []string{"name", "namespace"})

	gaugeVecs["pods"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "pod_status",
			Help:      "status of pod reported by kubernates",
		}, []string{"name", "namespace", "podPhase", "hostIP", "podIP", "reason", "message", "containerCount"})

	return gaugeVecs
}
