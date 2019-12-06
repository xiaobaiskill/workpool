package proxypool

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	proxypoolnums *prometheus.GaugeVec

	registry *prometheus.Registry
}

func (m *Metrics) ProxypoolnumInc(name string){
	m.proxypoolnums.WithLabelValues(name).Inc()
}

func (m *Metrics) ProxypoolnumDec(name string){
	m.proxypoolnums.WithLabelValues(name).Dec()
}

func (m *Metrics) GinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}).ServeHTTP(c.Writer,c.Request)
	}
}

func NewMetrics(metricsPrefix string)*Metrics {
	m := &Metrics{
		proxypoolnums: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: metricsPrefix,
			Name:      "proxypoolnums",
			Help:      "这里记录着各个代理池的个数",
		}, []string{"execting"}),
	}
	m.registry = prometheus.NewRegistry()
	m.registry.MustRegister(m.proxypoolnums)
	return m
}
