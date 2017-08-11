package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	beta "k8s.io/client-go/pkg/apis/apps/v1beta1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

// Data kubernate rest data
type Data struct {
}

// Discovery kubernate discovery client
type Discovery struct {
	client kubernetes.Interface
	role   string
}

// GathData kubernates data
type GathData struct {
	pods         *v1.PodList
	nodes        *v1.NodeList
	services     *v1.ServiceList
	endpoints    *v1.EndpointsList
	deployments  *v1beta1.DeploymentList
	daemonsets   *v1beta1.DaemonSetList
	statefulsets *beta.StatefulSetList
	stacks       map[Stack]*[]v1beta1.Deployment
	components   map[KubeComponent]*[]v1.Pod
}

// Run fetch kubernates data
func (d *Discovery) Run() *GathData {

	daemonsets, err := d.client.ExtensionsV1beta1().DaemonSets(api.NamespaceAll).List(v1.ListOptions{})
	if err != nil {
		fmt.Println(err, "error")
	}

	deployments, err := d.client.ExtensionsV1beta1().Deployments(api.NamespaceAll).List(v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	statefuls, err := d.client.AppsV1beta1().StatefulSets(api.NamespaceAll).List(v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	pods, err := d.client.Core().Pods(api.NamespaceAll).List(v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	nodes, err := d.client.Core().Nodes().List(v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	services, err := d.client.Core().Services(api.NamespaceAll).List(v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	endpoints, err := d.client.Core().Endpoints(api.NamespaceAll).List(v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	stacks := map[Stack]*[]v1beta1.Deployment{}

	for _, deployment := range deployments.Items {
		stack := Stack{
			Name:      stackName(deployment),
			Namespace: deployment.Namespace,
		}

		if deployments, ok := stacks[stack]; ok {
			*deployments = append(*deployments, deployment)
			stacks[stack] = deployments
		} else {
			stacks[stack] = &[]v1beta1.Deployment{deployment}
		}
	}

	components := map[KubeComponent]*[]v1.Pod{}
	for _, pod := range pods.Items {
		if pod.Namespace == "kube-system" {
			for _, component := range kubeComponents {
				if strings.Contains(pod.Name, component) {
					cmp := KubeComponent{
						Name:      component,
						Namespace: pod.Namespace,
					}

					if pods, ok := components[cmp]; ok {
						*pods = append(*pods, pod)
						components[cmp] = pods
					} else {
						components[cmp] = &[]v1.Pod{pod}
					}
				}
			}
		}
	}

	return &GathData{
		pods:         pods,
		nodes:        nodes,
		services:     services,
		endpoints:    endpoints,
		deployments:  deployments,
		stacks:       stacks,
		daemonsets:   daemonsets,
		statefulsets: statefuls,
		components:   components,
	}

}

// New new discovery instance
func (e *Exporter) New() (*Discovery, error) {

	var (
		kcfg *rest.Config
	)

	kcfg = &rest.Config{
		Host: e.APIServer.String(),
		TLSClientConfig: rest.TLSClientConfig{
			CAFile:   e.TLSConfig.CAFile,
			CertFile: e.TLSConfig.CertFile,
			KeyFile:  e.TLSConfig.KeyFile,
		},
		Insecure: e.TLSConfig.InsecureSkipVerify,
	}

	token := e.BearerToken

	if e.BearerTokenFile != "" {
		bf, err1 := ioutil.ReadFile(e.BearerTokenFile)
		if err1 != nil {
			return nil, err1
		}
		token = string(bf)
	}

	kcfg.BearerToken = token

	kcfg.UserAgent = "prometheus/kubernates-exporter"

	c, err2 := kubernetes.NewForConfig(kcfg)
	if err2 != nil {
		return nil, err2
	}

	return &Discovery{
		client: c,
	}, nil

}

func (e *Exporter) gatherData(ch chan<- prometheus.Metric) (*GathData, error) {

	discovery, err := e.New()
	if err != nil {
		fmt.Println(0, err)
		return nil, err
	}

	data := discovery.Run()
	return data, nil

}

// Stack group of deployment
type Stack struct {
	Name      string
	Namespace string
}

// KubeComponent group of pod
type KubeComponent struct {
	Name      string
	Namespace string
}

func stackName(deployment v1beta1.Deployment) string {
	return strings.Split(deployment.Name, "-")[0]
}

var kubeComponents = []string{"etcd", "kube-apiserver", "kube-controller-manager", "kube-dns", "kube-flannel", "kube-proxy", "kube-scheduler", "kubernetes-dashboard"}
