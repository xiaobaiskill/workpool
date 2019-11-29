package pool

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var metric *metrics

type metrics struct{
	jobTotalCounter prometheus.Counter
	executingGauge *prometheus.GaugeVec
}

func initMetrics(metricsPrefix string){
	metric = &metrics{
		jobTotalCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace:metricsPrefix,
			Name:"Job_total",
			Help:"This is used to calculate the total number of jobs",
		}),
		executingGauge:  prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace:metricsPrefix,
			Name:"Job_execting",
			Help:"This is used to measure the executing job",
		},[]string{"execting"}),
	}
}


func (m *metrics) jobTotalInc(){
	m.jobTotalCounter.Inc()
}

func (m *metrics) jobexectingInc(){
	m.executingGauge.WithLabelValues("job_execting").Inc()
}

func (m *metrics) jobexectingDec(){
	m.executingGauge.WithLabelValues("job_execting").Dec()
}



func (m *metrics) ginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer,c.Request)
	}
}
