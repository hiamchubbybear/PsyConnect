


package stats

import "google.golang.org/grpc/stats"


type MetricsRecorder interface {
	
	
	RecordInt64Count(handle *Int64CountHandle, incr int64, labels ...string)
	
	
	RecordFloat64Count(handle *Float64CountHandle, incr float64, labels ...string)
	
	
	RecordInt64Histo(handle *Int64HistoHandle, incr int64, labels ...string)
	
	
	RecordFloat64Histo(handle *Float64HistoHandle, incr float64, labels ...string)
	
	
	RecordInt64Gauge(handle *Int64GaugeHandle, incr int64, labels ...string)
}



type Metrics = stats.MetricSet


type Metric = string



func NewMetrics(metrics ...Metric) *Metrics {
	return stats.NewMetricSet(metrics...)
}
