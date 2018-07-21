package backend

// InitPrometheus : Set some channels to notify other theads when using Prometheus
import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// InitPrometheus : initialize prometheus
func (backend *Config) InitPrometheus(channel *chan bool, doneChannel *chan bool, promMetrics *chan Point) error {
	backend.channel = channel
	backend.doneChannel = doneChannel
	backend.promMetrics = promMetrics
	return nil
}

// Describe : Implementation of Prometheus Collector.Describe
func (backend *Config) Describe(ch chan<- *prometheus.Desc) {
	prometheus.NewGauge(prometheus.GaugeOpts{Name: "Dummy", Help: "Dummy"}).Describe(ch)
}

// Collect : Implementation of Prometheus Collector.Collect
func (backend *Config) Collect(ch chan<- prometheus.Metric) {

	stdlog.Println("Requested Metrics!")

	*backend.channel <- true

	for {
		select {
		case point := <-*backend.promMetrics:
			tags := point.GetTags(backend.NoArray, ",")
			labelNames := make([]string, len(tags))
			labelValues := make([]string, len(tags))
			i := 0
			for key, value := range tags {
				labelNames[i] = key
				labelValues[i] = value
				i++
			}
			key := fmt.Sprintf("vsphere_%s_%s_%s", point.Group, point.Counter, point.Rollup)
			desc := prometheus.NewDesc(key, "vSphere collected metric", labelNames, nil)
			metric, err := prometheus.NewConstMetric(desc, prometheus.GaugeValue, float64(point.Value), labelValues...)
			if err != nil {
				errlog.Println("Error creating prometheus metric")
			}
			ch <- metric
		case <-*backend.doneChannel:
			return
		}
	}
}