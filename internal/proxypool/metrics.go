package proxypool

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	proxypoolnums *prometheus.GaugeVec

	registry *prometheus.Registry
}

func (m *metrics) proxypoolnumset(name string, num float64){
	m.proxypoolnums.WithLabelValues(name).Set(num)
}

func (m *metrics) GinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}).ServeHTTP(c.Writer,c.Request)
	}
}

func NewMetrics(metricsPrefix string)*metrics {
	m := &metrics{
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
