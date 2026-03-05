

package stats

import (
	"fmt"

	estats "google.golang.org/grpc/experimental/stats"
	"google.golang.org/grpc/stats"
)





type MetricsRecorderList struct {
	
	metricsRecorders []estats.MetricsRecorder
}





func NewMetricsRecorderList(shs []stats.Handler) *MetricsRecorderList {
	var mrs []estats.MetricsRecorder
	for _, sh := range shs {
		if mr, ok := sh.(estats.MetricsRecorder); ok {
			mrs = append(mrs, mr)
		}
	}
	return &MetricsRecorderList{
		metricsRecorders: mrs,
	}
}

func verifyLabels(desc *estats.MetricDescriptor, labelsRecv ...string) {
	if got, want := len(labelsRecv), len(desc.Labels)+len(desc.OptionalLabels); got != want {
		panic(fmt.Sprintf("Received %d labels in call to record metric %q, but expected %d.", got, desc.Name, want))
	}
}



func (l *MetricsRecorderList) RecordInt64Count(handle *estats.Int64CountHandle, incr int64, labels ...string) {
	verifyLabels(handle.Descriptor(), labels...)

	for _, metricRecorder := range l.metricsRecorders {
		metricRecorder.RecordInt64Count(handle, incr, labels...)
	}
}



func (l *MetricsRecorderList) RecordFloat64Count(handle *estats.Float64CountHandle, incr float64, labels ...string) {
	verifyLabels(handle.Descriptor(), labels...)

	for _, metricRecorder := range l.metricsRecorders {
		metricRecorder.RecordFloat64Count(handle, incr, labels...)
	}
}



func (l *MetricsRecorderList) RecordInt64Histo(handle *estats.Int64HistoHandle, incr int64, labels ...string) {
	verifyLabels(handle.Descriptor(), labels...)

	for _, metricRecorder := range l.metricsRecorders {
		metricRecorder.RecordInt64Histo(handle, incr, labels...)
	}
}



func (l *MetricsRecorderList) RecordFloat64Histo(handle *estats.Float64HistoHandle, incr float64, labels ...string) {
	verifyLabels(handle.Descriptor(), labels...)

	for _, metricRecorder := range l.metricsRecorders {
		metricRecorder.RecordFloat64Histo(handle, incr, labels...)
	}
}



func (l *MetricsRecorderList) RecordInt64Gauge(handle *estats.Int64GaugeHandle, incr int64, labels ...string) {
	verifyLabels(handle.Descriptor(), labels...)

	for _, metricRecorder := range l.metricsRecorders {
		metricRecorder.RecordInt64Gauge(handle, incr, labels...)
	}
}
