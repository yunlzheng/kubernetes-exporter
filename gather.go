package main

import (
	"fmt"
	"io/ioutil"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
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
	pods        *v1.PodList
	nodes       *v1.NodeList
	services    *v1.ServiceList
	endpoints   *v1.EndpointsList
	deployments *v1beta1.DeploymentList
}

// Run fetch kubernates data
func (d *Discovery) Run() *GathData {

	fmt.Println("Discovery Run")

	deployments, err := d.client.ExtensionsV1beta1().Deployments(api.NamespaceAll).List(v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(deployments, "deployments")

	//pods := &v1.PodList{}
	pods, err := d.client.Core().Pods(api.NamespaceAll).List(v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(pods, "pods")

	nodes, err := d.client.Core().Nodes().List(v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(nodes, "nodes")

	services, err := d.client.Core().Services(api.NamespaceAll).List(v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println(services, "services")

	endpoints, err := d.client.Core().Endpoints(api.NamespaceAll).List(v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println(endpoints, "endpoints")

	return &GathData{
		pods:        pods,
		nodes:       nodes,
		services:    services,
		endpoints:   endpoints,
		deployments: deployments,
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
