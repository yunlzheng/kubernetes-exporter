package main

import (
	"net/url"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Exporter Sets up all the runtime and metrics
type Exporter struct {
	APIServer       *url.URL
	BasicAuth       *BasicAuth
	BearerToken     string
	BearerTokenFile string
	TLSConfig       *TLSConfig
	mutex           sync.RWMutex
	gaugeVecs       map[string]*prometheus.GaugeVec
}

func newExporter(APIServer *url.URL, BearerTokenFile string, TLSConfig *TLSConfig) *Exporter {
	gaugeVecs := addMetrics()
	return &Exporter{
		APIServer:       APIServer,
		BearerTokenFile: BearerTokenFile,
		TLSConfig:       TLSConfig,
		gaugeVecs:       gaugeVecs,
	}
}

type TLSConfig struct {
	CAFile             string
	CertFile           string
	KeyFile            string
	ServerName         string
	InsecureSkipVerify bool
	XXX                map[string]interface{}
}

type BasicAuth struct {
	Username string
	Password string
}

type KubernetesRole string

// The valid options for KubernetesRole.
const (
	KubernetesRoleNode     = "node"
	KubernetesRolePod      = "pod"
	KubernetesRoleService  = "service"
	KubernetesRoleEndpoint = "endpoints"
)
