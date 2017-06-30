package main

import (
	"fmt"
	"strconv"

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
				containerSize := len(pod.Status.ContainerStatuses)
				for _, item := range pod.Status.ContainerStatuses {
					if !item.Ready {
						state = 0
					}
				}
				e.gaugeVecs["pods"].With(prometheus.Labels{"name": pod.Name, "namespace": pod.Namespace, "podPhase": string(pod.Status.Phase), "hostIP": pod.Status.HostIP, "podIP": pod.Status.PodIP, "reason": pod.Status.Reason, "message": pod.Status.Message, "containerCount": strconv.Itoa(containerSize)}).Set(state)
			}

			for _, node := range gathData.nodes.Items {
				var state float64 = 1
				var nodeState = "Ready"
				var hostname = ""
				var externalIp = ""
				var internalIp = ""

				for _, address := range node.Status.Addresses {
					if address.Type == v1.NodeHostName {
						hostname = address.Address
					} else if address.Type == v1.NodeExternalIP {
						externalIp = address.Address
					} else if address.Type == v1.NodeInternalIP {
						internalIp = address.Address
					}
				}

				for _, item := range node.Status.Conditions {
					if item.Type != v1.NodeReady && item.Status == v1.ConditionTrue {
						state = 0
						nodeState = string(item.Type)
					}
				}

				e.gaugeVecs["nodes"].With(
					prometheus.Labels{
						"name":                    node.Name,
						"nodeState":               nodeState,
						"osImage":                 node.Status.NodeInfo.OSImage,
						"containerRuntimeVersion": node.Status.NodeInfo.ContainerRuntimeVersion,
						"kubeletVersion":          node.Status.NodeInfo.KubeletVersion,
						"operatingSystem":         node.Status.NodeInfo.OperatingSystem,
						"architecture":            node.Status.NodeInfo.Architecture,
						"hostname":                hostname,
						"externalIp":              externalIp,
						"internalIp":              internalIp,
					}).Set(state)
			}

			for _, deployment := range gathData.deployments.Items {
				var state float64 = 1
				for _, condition := range deployment.Status.Conditions {
					fmt.Println(condition)
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
