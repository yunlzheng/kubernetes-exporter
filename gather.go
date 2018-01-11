package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mojo-zd/go-library/traverse"
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

type ResourceStatus struct {
	resources map[string]bool
	status    chan string
}

// Run fetch kubernates data
func (d *Discovery) Run() (gathData *GathData) {
	over := false
	resourceStatus := ResourceStatus{
		resources: map[string]bool{"daemonSets": true, "deployments": true, "statefuls": true, "pods": true, "nodes": true, "services": true, "endpoints": true},
		status:    make(chan string),
	}

	gathData = &GathData{}
	go func() {
		daemonSets, err := d.client.ExtensionsV1beta1().DaemonSets(api.NamespaceAll).List(v1.ListOptions{})
		if err != nil {
			fmt.Printf("daemonSets error info is %s\n", err.Error())
		}
		gathData.daemonsets = daemonSets
		resourceStatus.status <- "daemonSets"
	}()

	go func() {
		deployments, err := d.client.ExtensionsV1beta1().Deployments(api.NamespaceAll).List(v1.ListOptions{})
		if err != nil {
			fmt.Printf("deployments error info is %s\n", err.Error())
		}
		gathData.deployments = deployments
		gathData.stacks = assembleStack(deployments)
		resourceStatus.status <- "deployments"
	}()

	go func() {
		statefuls, err := d.client.AppsV1beta1().StatefulSets(api.NamespaceAll).List(v1.ListOptions{})
		if err != nil {
			fmt.Printf("statefuls error info is %s\n", err.Error())
		}
		gathData.statefulsets = statefuls
		resourceStatus.status <- "statefuls"
	}()

	go func() {
		pods, err := d.client.CoreV1().Pods(api.NamespaceAll).List(v1.ListOptions{})
		if err != nil {
			fmt.Printf("pods error info is %s\n", err.Error())
		}
		gathData.pods = pods
		gathData.components = assembleComponent(pods)
		resourceStatus.status <- "pods"
	}()

	go func() {
		nodes, err := d.client.CoreV1().Nodes().List(v1.ListOptions{})
		if err != nil {
			fmt.Printf("nodes error info is %s\n", err.Error())
		}
		gathData.nodes = nodes
		resourceStatus.status <- "nodes"
	}()

	go func() {
		services, err := d.client.CoreV1().Services(api.NamespaceAll).List(v1.ListOptions{})
		if err != nil {
			fmt.Printf("services error info is %s\n", err.Error())
		}
		gathData.services = services
		resourceStatus.status <- "services"
	}()

	go func() {
		endpoints, err := d.client.CoreV1().Endpoints(api.NamespaceAll).List(v1.ListOptions{})
		if err != nil {
			fmt.Printf("endpoints error info is %s\n", err.Error())
		}
		gathData.endpoints = endpoints
		resourceStatus.status <- "endpoints"
	}()

	for {
		select {
		case state := <-resourceStatus.status:
			delete(resourceStatus.resources, state)
		default:
			if len(resourceStatus.resources) == 0 {
				over = true
			}
		}

		if over {
			break
		}
	}
	return
}

func assembleStack(deployments *v1beta1.DeploymentList) (stacks map[Stack]*[]v1beta1.Deployment) {
	stacks = map[Stack]*[]v1beta1.Deployment{}

	traverse.Iterator(deployments.Items, func(index int, value interface{}) (flag traverse.CYCLE_FLAG) {
		item := value.(v1beta1.Deployment)
		stack := Stack{
			Name:      stackName(item),
			Namespace: item.Namespace,
		}

		if deployment, ok := stacks[stack]; ok {
			*deployment = append(*deployment, item)
			stacks[stack] = deployment
		} else {
			stacks[stack] = &[]v1beta1.Deployment{item}
		}
		return
	})
	return
}

func assembleComponent(podList *v1.PodList) (components map[KubeComponent]*[]v1.Pod) {
	components = map[KubeComponent]*[]v1.Pod{}
	traverse.Iterator(podList.Items, func(index int, value interface{}) (flag traverse.CYCLE_FLAG) {
		pod := value.(v1.Pod)
		if pod.Namespace != "kube-system" {
			flag = traverse.CONTINUE_FLAT
			return
		}

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
		return
	})
	return
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
		fmt.Println(err)
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
