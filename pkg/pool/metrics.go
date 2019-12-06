package pool

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var metric *metrics

type metrics struct{
	jobTotalCounter prometheus.Counter
	executingGauge prometheus.Gauge

	registry *prometheus.Registry
}

func initMetrics(metricsPrefix string){
	metric = &metrics{
		jobTotalCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace:metricsPrefix,
			Name:"Job_total",
			Help:"总共运行了多少次job",
		}),
		executingGauge:  prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:metricsPrefix,
			Name:"Job_execting",
			Help:"正在运行中的job",
		}),
	}
	metric.registry = prometheus.NewRegistry()
	metric.registry.MustRegister(metric.jobTotalCounter,metric.executingGauge)
}


func (m *metrics) jobTotalInc(){
	m.jobTotalCounter.Inc()
}

func (m *metrics) jobexectingInc(){
	m.executingGauge.Inc()
}

func (m *metrics) jobexectingDec(){
	m.executingGauge.Dec()
}


func (m *metrics) ginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}).ServeHTTP(c.Writer,c.Request)
	}
}
