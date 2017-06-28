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

func newExporter(APIServer *url.URL, BearerToken string, BearerTokenFile string, TLSConfig *TLSConfig) *Exporter {
	gaugeVecs := addMetrics()
	return &Exporter{
		APIServer:       APIServer,
		BearerToken:     BearerToken,
		BearerTokenFile: BearerTokenFile,
		TLSConfig:       TLSConfig,
		gaugeVecs:       gaugeVecs,
	}
}

// TLSConfig kubernates client tls config
type TLSConfig struct {
	CAFile             string
	CertFile           string
	KeyFile            string
	ServerName         string
	InsecureSkipVerify bool
	XXX                map[string]interface{}
}

// BasicAuth kubernates client basic auth
type BasicAuth struct {
	Username string
	Password string
}
