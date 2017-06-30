package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func (e *Exporter) resetGaugeVecs() {
	for _, m := range e.gaugeVecs {
		m.Reset()
	}
}

// Describe describes all the metrics ever exported by the Rancher exporter
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.gaugeVecs {
		m.Describe(ch)
	}
}

// Collect function, called on by Prometheus Client library
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	e.resetGaugeVecs() // Clean starting point
	gathData, err := e.gatherData(ch)
	if err == nil {
		if gathData != nil {
			for _, pod := range gathData.pods.Items {

				var state float64 = 1
				var containerSize int = len(pod.Status.ContainerStatuses)
				for _, containerStatus := range pod.Status.ContainerStatuses {
					if !containerStatus.Ready {
						state = 0
					}
				}
				e.gaugeVecs["pods"].With(prometheus.Labels{"name": pod.Name, "namespace": pod.Namespace, "podPhase": string(pod.Status.Phase), "hostIP": pod.Status.HostIP, "podIP": pod.Status.PodIP, "reason": pod.Status.Reason, "message": pod.Status.Message, "containerSize": string(containerSize)}).Set(state)
			}

			for _, node := range gathData.nodes.Items {
				var state float64 = 1
				if node.Status.Phase != v1.NodeRunning {
					state = 0
				}
				e.gaugeVecs["nodes"].With(prometheus.Labels{"name": node.Name, "namespace": node.Namespace}).Set(state)
			}

			for _, deployment := range gathData.deployments.Items {
				var state float64 = 1
				for _, condition := range deployment.Status.Conditions {
					if condition.Type != v1beta1.DeploymentAvailable {
						state = 0
					}
				}
				e.gaugeVecs["deployments"].With(prometheus.Labels{"name": deployment.Name, "namespace": deployment.Namespace}).Set(state)
			}

		}
	}

	for _, m := range e.gaugeVecs {
		m.Collect(ch)
	}

}
