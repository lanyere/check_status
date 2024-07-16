package check_status

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promclient "github.com/prometheus/client_model/go"
	"log"
	"net/http"
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(errorCount)
	prometheus.MustRegister(transactionsProcessed)
}

// NewPoller initializes new poller with Config
func NewPoller(conf *Config) Poller {
	p := &poll{
		config: conf,
	}
	return p
}

// StopPolling stops polling and closes poll data
func (p *poll) StopPolling() {
	p.status = false
}

// StartPolling initializes ticker to call each provider per interval
func (p *poll) StartPolling() {
	if err := p.verifyOptions(); err != nil {
		return
	}

	p.status = true
	go p.startTicker()
}

// GetTransactionStatus helps you to get Status quicker, instead of going through database
func (p *poll) GetTransactionStatus(id string) (Status, error) {
	val, ok := p.data.Load(id)
	if !ok {
		return "", fmt.Errorf("transaction %s not found", id)
	}

	status, ok := val.(Status)
	if !ok {
		// which is almost impossible, but still, let's be safer
		return "", fmt.Errorf("transaction %s is not a string, value: %v | value type: %T", id, val, val)
	}

	return status, nil
}

// ExposeMetrics will expose all the metrics in the specified endpoint and port
func ExposeMetrics(endpoint, port string) {
	http.Handle(endpoint, promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// CollectMetrics if you want to handle them however you want to
func CollectMetrics() ([]*promclient.MetricFamily, error) {
	registry := prometheus.DefaultRegisterer.(*prometheus.Registry)
	return registry.Gather()
}
