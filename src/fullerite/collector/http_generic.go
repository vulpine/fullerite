package collector

import (
	"fullerite/metric"

	"net/http"
)

type errorHandler func(error)
type responseHandler func(*http.Response) []metric.Metric

type baseHttpCollector struct {
	baseCollector

	rspHandler responseHandler
	errHandler errorHandler

	endpoint string
}

// Collect first queries the config'd endpoint and then passes the results to the handler functions
func (base baseHttpCollector) Collect() {
	base.log.Info("Starting to collect metrics from ", base.endpoint)

	metrics := base.makeRequest()
	if metrics != nil {
		for _, m := range metrics {
			base.Channel() <- m
		}

		base.log.Info("Collected and sent ", len(metrics), " metrics")
	} else {
		base.log.Info("Sent no metrics because we didn't get any from the response")
	}
}

// makeRequest is what is responsible for actually doing the HTTP GET
func (base baseHttpCollector) makeRequest() []metric.Metric {
	if base.endpoint == "" {
		base.log.Warn("Ignoring attempt to make request because no endpoint provided")
	}

	rsp, err := http.Get(base.endpoint)
	if err != nil {
		base.errHandler(err)
		return nil
	}

	return base.rspHandler(rsp)
}
