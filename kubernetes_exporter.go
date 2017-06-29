package main

import (
	"flag"
	"net/http"
	"net/url"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/yunlzheng/kubernates-exporter/measure"
)

const (
	namespace = "kubernetes" // Used to prepand Prometheus metrics created by this exporter.
)

var (
	metricsPath     = getEnv("METRICS_PATH", "/metrics")
	listenAddress   = getEnv("LISTEN_ADDRESS", ":9174")
	endpoints       = []string{"nodes", "pods", "deployments"}
	roles           = []string{"node", "pod", "service", "endpoints"}
	logLevel        = getEnv("LOG_LEVEL", "info") // Optional - Set the logging level
	apiServer       = "https://kubernetes.default.svc"
	bearerToken     = ""
	bearerTokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	tslConfig       = &TLSConfig{
		CAFile: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
	}
)

func main() {

	flag.Parse()

	setLogLevel(logLevel)

	log.Info("Starting Kubernates Exporter for Kubernates")
	log.Info("Runtime Configuration")

	measure.Init()

	APIServer, err := url.Parse(apiServer)
	if err != nil {
		log.Printf("Kubernates APIServer invalidate")
		os.Exit(1)
	}

	Exporter := newExporter(APIServer, bearerToken, bearerTokenFile, tslConfig)

	prometheus.MustRegister(Exporter)

	// Setup HTTP handler
	http.Handle(metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		                <head><title>Kubernates exporter</title></head>
		                <body>
		                   <h1>kubernates exporter</h1>
		                   <p><a href='` + metricsPath + `'>Metrics</a></p>
		                   </body>
		                </html>
		              `))
	})
	log.Printf("Starting Server on port %s and path %s", listenAddress, metricsPath)
	log.Fatal(http.ListenAndServe(listenAddress, nil))

}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
