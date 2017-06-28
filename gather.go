package main

import (
	"fmt"
	"io/ioutil"

	// metaapi "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/rest"
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
	pods      *v1.PodList
	nodes     *v1.NodeList
	services  *v1.ServiceList
	endpoints *v1.EndpointsList
}

// Run fetch kubernates data
func (d *Discovery) Run() *GathData {
	rclient := d.client.Core().GetRESTClient()

	fmt.Println(rclient.APIVersion(), "api version")

	pods := &v1.PodList{}
	err := rclient.Get().Namespace(api.NamespaceAll).Resource("pods").FieldsSelectorParam(nil).Do().Into(pods)
	if err != nil {
		return nil
	}

	fmt.Println(pods, "pods")

	nodes := &v1.NodeList{}
	err = rclient.Get().Namespace(api.NamespaceAll).Resource("nodes").FieldsSelectorParam(nil).Do().Into(nodes)
	if err != nil {
		return nil
	}

	fmt.Println(nodes, "nodes")

	services := &v1.ServiceList{}
	err = rclient.Get().Namespace(api.NamespaceAll).Resource("services").FieldsSelectorParam(nil).Do().Into(services)
	if err != nil {
		return nil
	}

	fmt.Println(services, "services")

	endpoints := &v1.EndpointsList{}
	err = rclient.Get().Namespace(api.NamespaceAll).Resource("endpoints").FieldsSelectorParam(nil).Do().Into(endpoints)
	if err != nil {
		return nil
	}

	fmt.Println(endpoints, "endpoints")

	return &GathData{
		pods:      pods,
		nodes:     nodes,
		services:  services,
		endpoints: endpoints,
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
