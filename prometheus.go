package main

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/pkg/api/v1"
	beta "k8s.io/client-go/pkg/apis/apps/v1beta1"
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
				if containerSize == 0 {
					state = 0
				}
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
				var hostname, externalIP, internalIP string

				for _, address := range node.Status.Addresses {
					if address.Type == v1.NodeHostName {
						hostname = address.Address
					} else if address.Type == v1.NodeExternalIP {
						externalIP = address.Address
					} else if address.Type == v1.NodeInternalIP {
						internalIP = address.Address
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
						"externalIp":              externalIP,
						"internalIp":              internalIP,
					}).Set(state)
			}

			for _, deployment := range gathData.deployments.Items {
				var state = getDeploymentState(deployment)
				e.gaugeVecs["controller"].With(prometheus.Labels{"name": deployment.Name, "namespace": deployment.Namespace, "type": "deployment"}).Set(state)
			}

			for _, daemonset := range gathData.daemonsets.Items {
				var state = getDaemonSetState(daemonset)
				e.gaugeVecs["controller"].With(prometheus.Labels{"name": daemonset.Name, "namespace": daemonset.Namespace, "type": "daemonset"}).Set(state)
			}

			for _, statefulset := range gathData.statefulsets.Items {
				var state = getStatefulSetState(statefulset)
				e.gaugeVecs["controller"].With(prometheus.Labels{"name": statefulset.Name, "namespace": statefulset.Namespace, "type": "statefulset"}).Set(state)
			}

			for stack, deployments := range gathData.stacks {
				var state float64 = 1
				for _, deployment := range *deployments {
					if getDeploymentState(deployment) == 0 {
						state = 0
						break
					}
				}
				e.gaugeVecs["stacks"].With(prometheus.Labels{"name": stack.Name, "namespace": stack.Namespace}).Set(state)
			}

			for component, pods := range gathData.components {
				var state float64 = 1
				for _, pod := range *pods {
					for _, item := range pod.Status.ContainerStatuses {
						if !item.Ready {
							state = 0
						}
					}
				}
				e.gaugeVecs["components"].With(prometheus.Labels{"name": component.Name, "namespace": component.Namespace}).Set(state)
			}

		}
	}

	for _, m := range e.gaugeVecs {
		m.Collect(ch)
	}

}

func getDaemonSetState(daemonset v1beta1.DaemonSet) float64 {
	var state float64 = 1
	if daemonset.Status.NumberReady != daemonset.Status.DesiredNumberScheduled {
		state = 0
	}
	return state
}

func getDeploymentState(deployment v1beta1.Deployment) float64 {
	var state float64 = 1
	if deployment.Status.Replicas != deployment.Status.AvailableReplicas {
		state = 0
	}
	return state
}

func getStatefulSetState(stateful beta.StatefulSet) float64 {
	var state float64 = 1
	if *stateful.Spec.Replicas  != stateful.Status.Replicas {
		state = 0
	}
	return state
}
