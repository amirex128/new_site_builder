package exporter

import "net/http"

type BasePrometheusExporter struct {
}

func (BasePrometheusExporter) Handler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
