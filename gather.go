package main

import (
	"fmt"
	"io/ioutil"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/pkg/api"
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

// Run fetch kubernates data
func (d *Discovery) Run() {
	rclient := d.client.Core().GetRESTClient()
	fmt.Println(rclient.APIVersion(), "api version")
	pods := rclient.Get().Namespace(api.NamespaceAll).Resource("pods").FieldsSelectorParam(nil).Do()
	nodes := rclient.Get().Namespace(api.NamespaceAll).Resource("nodes").FieldsSelectorParam(nil).Do()
	services := rclient.Get().Namespace(api.NamespaceAll).Resource("services").FieldsSelectorParam(nil).Do()
	endpoints := rclient.Get().Namespace(api.NamespaceAll).Resource("endpoints").FieldsSelectorParam(nil).Do()
	fmt.Println(pods.Raw())
	fmt.Println(nodes.Raw())
	fmt.Println(services.Raw())
	fmt.Println(endpoints.Raw())
}

func (d *Discovery) getNamespaces() []string {
	return []string{api.NamespaceAll}
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

func (e *Exporter) gatherData(ch chan<- prometheus.Metric) (*Data, error) {

	fmt.Println("New Discovery")

	discovery, err := e.New()
	if err != nil {
		fmt.Println(0, err)
		return nil, err
	}

	fmt.Println(" Discovery Run")
	discovery.Run()

	var data = new(Data)
	return data, nil

}
